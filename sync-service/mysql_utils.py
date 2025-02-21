import mysql.connector
import logging
from config import MYSQL_HOST, MYSQL_PORT, MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - [%(funcName)s:%(lineno)d] - %(message)s'
)

def get_patient_details(patient_id):
    """
    Retrieves patient details from the OpenMRS database using patient_id.

    Args:
        patient_id (int): The ID of the patient.

    Returns:
        dict: A dictionary containing the patient's details, or None if not found.
    """
    try:
        connection = mysql.connector.connect(
            host=MYSQL_HOST,
            port=MYSQL_PORT,
            user=MYSQL_USER,
            password=MYSQL_PASSWORD,
            database=MYSQL_DATABASE
        )
        cursor = connection.cursor(dictionary=True)

        query = """
            SELECT
                p.patient_id,
                pn.given_name,
                pn.family_name,
                pa.address1 AS street,
                pa.city_village AS city,
                ppn.value AS phone,
                pea.value AS email,
                per.birthdate,
                per.gender,
                '1' AS customer_rank,
                'person' AS company_type,
                MAX(CASE WHEN pi.identifier_type = 3 THEN pi.identifier END) AS ref
            FROM patient p
            LEFT JOIN person per ON per.person_id = p.patient_id
            LEFT JOIN person_name pn ON p.patient_id = pn.person_id AND pn.preferred = 1
            LEFT JOIN person_address pa ON p.patient_id = pa.person_id AND pa.preferred = 1
            LEFT JOIN person_attribute ppn ON p.patient_id = ppn.person_id
            LEFT JOIN patient_identifier pi ON p.patient_id = pi.patient_id
            AND ppn.person_attribute_type_id = (
                SELECT
                    person_attribute_type_id
                FROM
                    person_attribute_type
                WHERE
                    name = 'phone_number'
            )
            LEFT JOIN person_attribute pea ON p.patient_id = pea.person_id
            AND pea.person_attribute_type_id = (
                SELECT
                    person_attribute_type_id
                FROM
                    person_attribute_type
                WHERE
                    name = 'email'
            )
            WHERE p.patient_id = %s;
        """

        cursor.execute(query, (patient_id,))
        result = cursor.fetchone()

        if result:
            logging.info(f"Patient details retrieved successfully for patient_id: {patient_id}")
            return result
        else:
            logging.warning(f"Patient not found with patient_id: {patient_id}")
            return None

    except mysql.connector.Error as error:
        logging.error(f"Error fetching patient details from MySQL: {error}")
        return None
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()
            logging.info("MySQL connection closed")
