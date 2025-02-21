import requests
import logging
from typing import Dict, Optional, Union
from config import ODOO_URL, ODOO_DB, ODOO_USER, ODOO_PASSWORD
import json

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - [%(funcName)s:%(lineno)d] - %(message)s'
)

def create_customer(
    name: str,
    phone: Optional[str] = None,
    email: Optional[str] = None,
    street: Optional[str] = None,
    city: Optional[str] = None,
    is_company: bool = False
) -> Dict[str, Union[str, int]]:
    url = ODOO_URL
    db = ODOO_DB
    username = ODOO_USER
    password = ODOO_PASSWORD

    if not url.endswith('/'):
        url += '/'

    # Prepare customer data
    customer_data = {
        'name': name,
        'customer_rank': 1,
        'company_type': 'company' if is_company else 'person'
    }

    # Add optional fields if provided
    if phone:
        customer_data['phone'] = phone
    if email:
        customer_data['email'] = email
    if street:
        customer_data['street'] = street
    if city:
        customer_data['city'] = city

    # 1. Login to get session id
    login_endpoint = f"{url}web/session/authenticate"
    login_payload = {
        "jsonrpc": "2.0",
        "method": "call",
        "params": {
            "db": db,
            "login": username,
            "password": password
        }
    }

    try:
        session = requests.Session()
        login_response = session.post(
            login_endpoint,
            json=login_payload,
            timeout=10
        )
        login_response.raise_for_status()
        login_result = login_response.json()

        # 2. Create customer using the established session
        create_endpoint = f"{url}web/dataset/call_kw/res.partner/create"  # Changed to /jsonrpc endpoint
        create_payload = {
            "jsonrpc": "2.0",
            "method": "call",
            "params": {
                "model": "res.partner",
                "method": "create",
                "args": [customer_data],
                "kwargs": {}
            }
        }

        create_response = session.post(
            create_endpoint,
            json=create_payload,
            timeout=10,
            headers={'Content-Type': 'application/json'}  # Required header for jsonrpc
        )
        create_response.raise_for_status()

        create_result = create_response.json()

        if 'error' in create_result:
            logging.error(f"Create error: {create_result['error']['data']['message']}")  # Use logging
            return {
                'success': False,
                'error': create_result['error']['data']['message'],
                'code': 'odoo_create_error'
            }

        customer_id = create_result.get('result')  # Get ID from 'result'
        if customer_id is None:
            logging.error(f"Customer ID not found in create response: {create_result}")  # Use logging
            return {
                'success': False,
                'error': 'Customer ID not found in response',
                'code': 'odoo_id_error'
            }

        return {
            'success': True,
            'customer_id': customer_id,
            'message': 'Customer created successfully'
        }

    except requests.exceptions.RequestException as e:
        logging.exception("Request exception occurred.")  # Log the full traceback
        return {
            'success': False,
            'error': str(e),
            'code': 'request_error'
        }
    except json.JSONDecodeError:
        logging.exception("JSON decode error.")  # Log the full traceback
        return {
            'success': False,
            'error': 'Invalid JSON response from server',
            'code': 'json_error'
        }
    except Exception as e:
        logging.exception("Unexpected error occurred.")  # Log the full traceback
        return {
            'success': False,
            'error': f'Unexpected error: {str(e)}',
            'code': 'unknown_error'
        }
