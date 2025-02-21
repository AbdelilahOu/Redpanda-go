class PatientEvent:
    def __init__(self, schema, payload):
        self.schema = schema
        self.payload = payload

class Patient:
    def __init__(self, patient_id=None, **kwargs):
        # Using lowercase variable names to match the database fields
        self.patient_id = patient_id
        # Store any additional fields that might come from the database
        for key, value in kwargs.items():
            setattr(self, key, value)
