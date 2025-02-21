import os

# Kafka Configuration
KAFKA_BOOTSTRAP_SERVERS = os.environ.get("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
# KAFKA_GROUP_ID = os.environ.get("KAFKA_GROUP_ID", "patient")
# DEBEZIUM_CONNECTOR_HOST = os.environ.get("DEBEZIUM_CONNECTOR_HOST", "http://localhost:8083")
# DEBEZIUM_CONNECTOR_NAME = os.environ.get("DEBEZIUM_CONNECTOR_NAME", "patient")

# Odoo API Configuration
ODOO_URL = os.environ.get("ODOO_URL", "http://localhost:8069")
ODOO_DB = os.environ.get("ODOO_DB", "odoo")
ODOO_USER = os.environ.get("ODOO_USER", "admin")
ODOO_PASSWORD = os.environ.get("ODOO_PASSWORD", "admin")
ODOO_PATIENT_MODEL = os.environ.get("ODOO_PATIENT_MODEL", "res.partner")

# MySQL Configuration
MYSQL_HOST = os.environ.get("MYSQL_HOST", "localhost")
MYSQL_PORT = int(os.environ.get("MYSQL_PORT", "3306"))
MYSQL_USER = os.environ.get("MYSQL_USER", "openmrs-user")
MYSQL_PASSWORD = os.environ.get("MYSQL_PASSWORD", "password")
MYSQL_DATABASE = os.environ.get("MYSQL_DATABASE", "openmrs")

# DCM4CHEE Configuration
DCM4CHEE_URL = os.environ.get("DCM4CHEE_URL", "https://dev.modoock.com/dcm4chee-arc/aets/DCM4CHEE/rs/patients")
DCM4CHEE_AET = os.environ.get("DCM4CHEE_AET", "DCM4CHEE")
