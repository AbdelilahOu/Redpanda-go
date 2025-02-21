import asyncio
import logging
import argparse
from kafka_utils import create_kafka_topic, check_topic_details
from consumers import patient_consumer
from consumers import order_consumer

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - [%(funcName)s:%(lineno)d] - %(message)s'
)

async def main():
    """Main function to orchestrate the consumers and topic creation concurrently."""
    try:
        # Define coroutines for each consumer and topic setup
        async def patient_task():
            create_kafka_topic("patient.openmrs.patient")
            check_topic_details("patient.openmrs.patient")
            await patient_consumer.consume_messages("patient.openmrs.patient", "patient")

        async def order_task():
            create_kafka_topic("order.openmrs.order")
            check_topic_details("order.openmrs.order")
            await order_consumer.consume_messages("order.openmrs.order", "order")

        # Run the tasks concurrently using asyncio.gather
        await asyncio.gather(
            patient_task(),
            order_task(),
        )

    except Exception as e:
        logging.exception(f"Error in main function: {e}")

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logging.info("Consumer stopped by user")
    except Exception as e:
        logging.exception(f"Unhandled exception: {e}")
