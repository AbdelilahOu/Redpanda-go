import logging
from kafka import KafkaAdminClient, KafkaConsumer
from kafka.admin import NewTopic
from kafka.errors import TopicAlreadyExistsError, KafkaError
from config import KAFKA_BOOTSTRAP_SERVERS
import json

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - [%(funcName)s:%(lineno)d] - %(message)s'
)

def check_topic_details(topic_name):
    """Check Kafka topic details and configuration."""
    admin_client = KafkaAdminClient(bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS)
    try:
        topics = admin_client.list_topics()
        logging.info(f"Available topics: {topics}")

        if topic_name in topics:
            topic_details = admin_client.describe_topics([topic_name])
            logging.info(f"Topic details for '{topic_name}': {topic_details}")
            return True
        return False
    except Exception as e:
        logging.error(f"Error checking topics: {e}")
        return False
    finally:
        admin_client.close()

def create_kafka_topic(topic_name):
    """Creates a Kafka topic if it doesn't exist."""
    try:
        admin_client = KafkaAdminClient(bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS)

        # Check if topic exists
        if check_topic_details(topic_name):
            logging.info(f"Topic '{topic_name}' already exists")
            return

        topic_list = [
            NewTopic(
                name=topic_name,
                num_partitions=3,  # Adjust based on your needs
                replication_factor=1,
                topic_configs={
                    'cleanup.policy': 'delete',
                    'retention.ms': '604800000',  # 7 days
                }
            )
        ]
        admin_client.create_topics(new_topics=topic_list, validate_only=False)
        logging.info(f"Topic '{topic_name}' created successfully")
    except TopicAlreadyExistsError:
        logging.info(f"Topic '{topic_name}' already exists")
    except Exception as e:
        logging.error(f"Error creating topic: {e}")
    finally:
        if 'admin_client' in locals():
            admin_client.close()


def create_kafka_consumer(topic, group_id):
    """Creates a Kafka consumer instance."""
    try:
        consumer = KafkaConsumer(
            topic,
            bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
            group_id=group_id,
            auto_offset_reset='earliest',
            enable_auto_commit=False,
            api_version=(0, 10, 1),
            value_deserializer=lambda x: json.loads(x.decode('utf-8')),
            session_timeout_ms=30000,
            heartbeat_interval_ms=10000,
        )
        return consumer
    except KafkaError as e:
        logging.error(f"Error creating Kafka consumer: {e}")
        return None
    except Exception as e:
        logging.error(f"Unexpected error creating Kafka consumer: {e}")
        return None
