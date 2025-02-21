import asyncio
import json
import logging
from kafka.errors import KafkaError
from kafka.structs import TopicPartition, OffsetAndMetadata

from kafka_utils import create_kafka_consumer
from odoo_utils import create_customer
from dcm4chee_utils import create_dcm4chee_patient
from mysql_utils import get_patient_details
from config import ODOO_URL, ODOO_DB, ODOO_USER, ODOO_PASSWORD

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - [%(funcName)s:%(lineno)d] - %(message)s'
)

async def handle_create(patient_data):
    """Handles the creation of a patient, retrieves additional details, and pushes to Odoo and DCM4CHEE."""
    if not patient_data:
        logging.warning("handle_create: Received empty patient data.")
        return

    patient_id = patient_data.get('patient_id')
    if not patient_id:
        logging.warning("handle_create: Patient ID not found in data.")
        return

    try:
        logging.info(f"handle_create: Received patient created event. Patient ID: {patient_id}")

        # Fetch patient details from OpenMRS
        patient_details = get_patient_details(patient_id)

        if not patient_details:
            logging.error(f"handle_create: Could not retrieve details for patient ID: {patient_id}")
            return

        # Prepare data for Odoo
        odoo_data = {
            'name': f"{patient_details.get('given_name', '')} {patient_details.get('family_name', '')}".strip() or "Unknown Patient",
            'customer_rank': patient_details.get('customer_rank'),
            'company_type': patient_details.get('company_type'),
            'phone': patient_details.get('phone'),
            'email': patient_details.get('email'),
            'street': patient_details.get('street'),
            'city': patient_details.get('city'),
            'ref': patient_details.get('ref'),
        }

        logging.info(f"handle_create: Data to be sent to Odoo: {odoo_data}")

        # Create Customer in Odoo
        odoo_result = await asyncio.to_thread(create_customer,
            name=odoo_data['name'],
            phone=odoo_data['phone'],
            email=odoo_data['email'],
            street=odoo_data['street'],
            city=odoo_data['city']
        )

        if odoo_result and odoo_result['success']:
            odoo_patient_id = odoo_result['customer_id']
            logging.info(f"handle_create: Patient created in Odoo with ID: {odoo_patient_id}")
        else:
            logging.error(f"handle_create: Failed to create patient in Odoo: {odoo_result}")
            return  # Stop if Odoo creation fails

        # Create Patient in DCM4CHEE
        dcm4chee_success = await create_dcm4chee_patient(patient_details)
        if dcm4chee_success:
            logging.info(f"handle_create: Patient created successfully in DCM4CHEE.")
        else:
            logging.error("handle_create: Failed to create patient in DCM4CHEE.")

    except Exception as e:
        logging.exception(f"handle_create: Error handling create operation: {e}")

async def handle_update(before_data, after_data):
    """Handles the update of a patient."""
    if not before_data or not after_data:
        return
    try:
        logging.info(f"Patient updated from: {before_data} to: {after_data}")
    except Exception as e:
        logging.error(f"Error handling update operation: {e}")

async def handle_delete(patient_data):
    """Handles the deletion of a patient."""
    if not patient_data:
        return
    try:
        logging.info(f"Patient deleted: {patient_data}")
    except Exception as e:
        logging.error(f"Error handling delete operation: {e}")

async def handle_read(patient_data):
    """Handles the read operation of a patient."""
    if not patient_data:
        return
    try:
        logging.info(f"Patient read: {patient_data}")
    except Exception as e:
        logging.error(f"Error handling read operation: {e}")

async def consume_messages(topic, group_id):
    """Consumes messages from the Kafka topic with proper error handling."""
    consumer = create_kafka_consumer(topic, group_id)

    if not consumer:
        logging.error("Failed to create Kafka consumer. Exiting.")
        return

    # Wait for partition assignment
    waiting_time = 0
    max_waiting_time = 60  # Maximum waiting time in seconds
    while not consumer.assignment() and waiting_time < max_waiting_time:
        consumer.poll(timeout_ms=1000)
        waiting_time += 1
        logging.info(f"Waiting for partition assignment... {waiting_time}s")

    if not consumer.assignment():
        logging.error("No partitions assigned after timeout")
        consumer.close()
        return

    logging.info(f"Assigned partitions: {consumer.assignment()}")

    try:
        for message in consumer:
            try:
                logging.info(f"Processing message from partition {message.partition} at offset {message.offset}")
                event = message.value

                if not isinstance(event, dict) or 'payload' not in event:
                    logging.warning(f"Skipping invalid message format: {event}")
                    tp = TopicPartition(message.topic, message.partition)
                    if tp in consumer.assignment():
                        consumer.commit({tp: OffsetAndMetadata(message.offset + 1, None)})
                    continue

                op = event['payload'].get('op')
                logging.info(f"Processing operation: {op}")

                if op == 'c':
                    await handle_create(event['payload'].get('after'))
                elif op == 'u':
                    await handle_update(event['payload'].get('before'), event['payload'].get('after'))
                elif op == 'd':
                    await handle_delete(event['payload'].get('before'))
                elif op == 'r':
                    await handle_read(event['payload'].get('after'))
                else:
                    logging.warning(f"Unknown operation type: {op}")

                # Commit offset after successful processing
                tp = TopicPartition(message.topic, message.partition)
                if tp in consumer.assignment():
                    try:
                        consumer.commit({tp: OffsetAndMetadata(message.offset + 1, None)})
                        logging.info(f"Successfully committed offset {message.offset + 1} for partition {message.partition}")
                    except Exception as commit_error:
                        logging.error(f"Error committing offset: {commit_error}")
                else:
                    logging.warning(f"Partition {tp} not assigned to this consumer")

            except json.JSONDecodeError as e:
                logging.error(f"Error decoding JSON: {e}, raw message: {message.value}")
                continue
            except Exception as e:
                logging.exception(f"Error processing message: {e}")
                continue

    except KafkaError as e:
        logging.error(f"Kafka consumer error: {e}")
    finally:
        logging.info("Closing consumer")
        consumer.close()
