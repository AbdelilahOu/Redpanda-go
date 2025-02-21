import asyncio
import requests
import logging
from datetime import datetime
from config import DCM4CHEE_URL

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - [%(funcName)s:%(lineno)d] - %(message)s'
)

def format_birth_date(birthdate):
    """Formats a date object to 'YYYYMMDD' string format."""
    if isinstance(birthdate, datetime):
        return birthdate.strftime('%Y%m%d')
    return None # Or raise an exception if birthdate is always expected to be a date

async def create_dcm4chee_patient(patient_details):
    """Creates a patient record in DCM4CHEE."""
    try:
        dcm_hl7_patient = {
            "00100010": {
                "vr": "PN",
                "Value": [
                    {
                        "Alphabetic": f"{patient_details.get('given_name', '')}^{patient_details.get('family_name', '')}",
                    },
                ],
            },
            "00100020": {"vr": "LO", "Value": [str(patient_details.get('patient_id'))]},
            "00100021": {"vr": "LO", "Value": ["ADT1"]},
            "00100030": {
                "vr": "DA",
                "Value": [format_birth_date(patient_details.get('birthdate'))],
            },
            "00100040": {"vr": "CS", "Value": [patient_details.get('gender')]},
        }

        response = await asyncio.to_thread(requests.post,
            DCM4CHEE_URL,
            json=dcm_hl7_patient,
            headers={"Content-Type": "application/json"}
        )
        response.raise_for_status()

        logging.info(f"DCM4CHEE patient creation successful. Status code: {response.status_code}")
        return True

    except requests.exceptions.RequestException as e:
        logging.error(f"Error creating DCM4CHEE patient: {e}")
        return False
    except Exception as e:
        logging.exception(f"Unexpected error creating DCM4CHEE patient: {e}")
        return False
