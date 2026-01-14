#!/bin/bash

# Quick API test - just the essentials
BASE_URL="https://ubiquitous-fishstick-697px5q7pv42rqjj-8080.app.github.dev/"

echo "🏥 Quick API Test"
echo "================="

# Health check
echo -e "\n1. Health Check..."
curl -s "$BASE_URL/health" | jq '.'

# Register & login
echo -e "\n2. Register Patient..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "quick@test.com",
    "password": "password123",
    "role": "patient",
    "first_name": "Quick",
    "last_name": "Test"
  }')

TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.tokens.access_token')

if [ "$TOKEN" != "null" ]; then
  echo "✓ Registered successfully"
  echo "Token: ${TOKEN:0:30}..."
  
  # Get profile
  echo -e "\n3. Get Profile..."
  curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/patient/profile" | jq '.'
  
  # Create study
  echo -e "\n4. Create Study..."
  curl -s -X POST "$BASE_URL/api/v1/studies" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "modality": "CBCT",
      "study_date": "2026-01-14"
    }' | jq '.'
  
  echo -e "\n✓ Quick test complete!"
else
  echo "✗ Registration failed"
  echo "$REGISTER_RESPONSE" | jq '.'
fi
