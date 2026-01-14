#!/bin/bash

# Helper script to activate clinic in database
CLINIC_ID=$1

if [ -z "$CLINIC_ID" ]; then
  echo "Usage: ./activate_clinic.sh <clinic_id>"
  exit 1
fi

echo "Activating clinic: $CLINIC_ID"

docker exec dental_postgres psql -U postgres -d dental_marketplace -c \
  "UPDATE clinics SET is_active = true WHERE id = '$CLINIC_ID';"

echo "Clinic activated!"
