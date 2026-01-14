#!/bin/bash

# ============================================
# Dental Marketplace API Test Script
# ============================================

# CONFIGURATION - Update this with your codespace URL
BASE_URL="https://ubiquitous-fishstick-697px5q7pv42rqjj-8080.app.github.dev"
# For GitHub Codespace, it will be something like:
# BASE_URL="https://your-codespace-name-8080.app.github.dev"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if jq is available
JQ_AVAILABLE=false
if command -v jq &> /dev/null; then
    JQ_AVAILABLE=true
fi

# Global variables to store tokens and IDs
PATIENT_TOKEN=""
CLINIC_TOKEN=""
PATIENT_ID=""
CLINIC_ID=""
STUDY_ID=""
PLAN_VERSION_ID=""
OFFER_REQUEST_ID=""
OFFER_ID=""
ORDER_ID=""

# Helper function to print section headers
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

# Helper function to print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Helper function to print error
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Helper function to print info
print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Helper function to pretty print JSON
pretty_json() {
    if [ "$JQ_AVAILABLE" = true ]; then
        echo "$1" | jq '.' 2>/dev/null || echo "$1"
    else
        echo "$1"
    fi
}

# Helper function to extract JSON value
json_value() {
    local json=$1
    local key=$2
    
    if [ "$JQ_AVAILABLE" = true ]; then
        echo "$json" | jq -r "$key" 2>/dev/null
    else
        # Fallback: simple grep/sed extraction
        echo "$json" | grep -o "\"$key\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" | sed 's/.*:"\(.*\)"/\1/'
    fi
}

# Helper function to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    
    if [ -n "$token" ]; then
        curl -s -X "$method" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "$data" \
            "$BASE_URL$endpoint"
    else
        curl -s -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint"
    fi
}

# ============================================
# TEST 1: Health Check
# ============================================
test_health() {
    print_header "TEST 1: Health Check"
    
    response=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$BASE_URL/health")
    http_code=$(echo "$response" | grep "HTTP_CODE" | cut -d: -f2)
    body=$(echo "$response" | sed '/HTTP_CODE/d')
    
    print_info "Response: $body"
    print_info "HTTP Code: $http_code"
    
    if [ "$http_code" = "200" ] || echo "$body" | grep -q "healthy"; then
        print_success "API is healthy"
        pretty_json "$body"
    else
        print_error "API health check failed (but continuing...)"
        echo "$body"
    fi
}

# ============================================
# TEST 2: Register Patient
# ============================================
test_register_patient() {
    print_header "TEST 2: Register Patient"
    
    data='{
        "email": "patient@test.com",
        "password": "password123",
        "role": "patient",
        "first_name": "John",
        "last_name": "Doe"
    }'
    
    response=$(api_call "POST" "/api/v1/auth/register" "$data")
    
    if echo "$response" | grep -q "access_token"; then
        print_success "Patient registered successfully"
        
        if [ "$JQ_AVAILABLE" = true ]; then
            PATIENT_TOKEN=$(echo "$response" | jq -r '.tokens.access_token')
            PATIENT_ID=$(echo "$response" | jq -r '.user.id')
        else
            PATIENT_TOKEN=$(echo "$response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
            PATIENT_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | head -1)
        fi
        
        print_info "Patient ID: $PATIENT_ID"
        print_info "Token: ${PATIENT_TOKEN:0:30}..."
        pretty_json "$response"
    else
        print_error "Patient registration failed"
        echo "$response"
    fi
}

# ============================================
# TEST 3: Login Patient
# ============================================
test_login_patient() {
    print_header "TEST 3: Login Patient"
    
    data='{
        "email": "patient@test.com",
        "password": "password123"
    }'
    
    response=$(api_call "POST" "/api/v1/auth/login" "$data")
    
    if echo "$response" | grep -q "access_token"; then
        print_success "Patient logged in successfully"
        
        if [ "$JQ_AVAILABLE" = true ]; then
            PATIENT_TOKEN=$(echo "$response" | jq -r '.tokens.access_token')
        else
            PATIENT_TOKEN=$(echo "$response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        fi
        
        pretty_json "$response"
    else
        print_error "Patient login failed"
        echo "$response"
    fi
}

# ============================================
# TEST 4: Get Patient Profile
# ============================================
test_get_patient_profile() {
    print_header "TEST 4: Get Patient Profile"
    
    response=$(curl -s -H "Authorization: Bearer $PATIENT_TOKEN" "$BASE_URL/api/v1/patient/profile")
    
    if echo "$response" | grep -q "first_name"; then
        print_success "Patient profile retrieved"
        pretty_json "$response"
    else
        print_error "Failed to get patient profile"
        echo "$response"
    fi
}

# ============================================
# TEST 5: Update Patient Profile
# ============================================
test_update_patient_profile() {
    print_header "TEST 5: Update Patient Profile"
    
    data='{
        "phone": "+1234567890",
        "preferred_city": "Moscow",
        "preferred_district": "Central",
        "preferred_price_segment": "business"
    }'
    
    response=$(api_call "PUT" "/api/v1/patient/profile" "$data" "$PATIENT_TOKEN")
    
    if echo "$response" | grep -q "Moscow"; then
        print_success "Patient profile updated"
        pretty_json "$response"
    else
        print_error "Failed to update patient profile"
        echo "$response"
    fi
}

# ============================================
# TEST 6: Create Study
# ============================================
test_create_study() {
    print_header "TEST 6: Create Study"
    
    data='{
        "modality": "CBCT",
        "study_date": "2026-01-14"
    }'
    
    response=$(api_call "POST" "/api/v1/studies" "$data" "$PATIENT_TOKEN")
    
    if echo "$response" | grep -q "id"; then
        print_success "Study created successfully"
        
        if [ "$JQ_AVAILABLE" = true ]; then
            STUDY_ID=$(echo "$response" | jq -r '.id')
        else
            STUDY_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | head -1)
        fi
        
        print_info "Study ID: $STUDY_ID"
        pretty_json "$response"
    else
        print_error "Failed to create study"
        echo "$response"
    fi
}

# ============================================
# TEST 7: Get Study
# ============================================
test_get_study() {
    print_header "TEST 7: Get Study"
    
    if [ -z "$STUDY_ID" ]; then
        print_error "No study ID available. Skipping."
        return
    fi
    
    response=$(curl -s -H "Authorization: Bearer $PATIENT_TOKEN" "$BASE_URL/api/v1/studies/$STUDY_ID")
    
    if echo "$response" | grep -q "$STUDY_ID"; then
        print_success "Study retrieved successfully"
        pretty_json "$response"
    else
        print_error "Failed to get study"
        echo "$response"
    fi
}

# ============================================
# TEST 8: Create Treatment Plan
# ============================================
test_create_plan() {
    print_header "TEST 8: Create Treatment Plan"
    
    if [ -z "$STUDY_ID" ]; then
        print_error "No study ID available. Skipping."
        return
    fi
    
    data='{
        "study_id": "'$STUDY_ID'",
        "source": "manual",
        "items": [
            {
                "tooth_number": 16,
                "specialty": "therapy",
                "procedure_code": "D2391",
                "procedure_name": "Resin-based composite - one surface",
                "diagnosis": "Caries",
                "quantity": 1
            },
            {
                "tooth_number": 26,
                "specialty": "orthopedics",
                "procedure_code": "D2740",
                "procedure_name": "Crown - porcelain/ceramic",
                "diagnosis": "Crown needed",
                "quantity": 1
            },
            {
                "tooth_number": 36,
                "specialty": "surgery",
                "procedure_code": "D7140",
                "procedure_name": "Extraction, erupted tooth",
                "diagnosis": "Non-restorable tooth",
                "quantity": 1
            }
        ]
    }'
    
    response=$(api_call "POST" "/api/v1/plans" "$data" "$PATIENT_TOKEN")
    
    if echo "$response" | grep -q "plan_items"; then
        print_success "Treatment plan created successfully"
        
        if [ "$JQ_AVAILABLE" = true ]; then
            PLAN_VERSION_ID=$(echo "$response" | jq -r '.id')
        else
            PLAN_VERSION_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | head -1)
        fi
        
        print_info "Plan Version ID: $PLAN_VERSION_ID"
        pretty_json "$response"
    else
        print_error "Failed to create treatment plan"
        echo "$response"
    fi
}

# ============================================
# TEST 9: Get Treatment Plan
# ============================================
test_get_plan() {
    print_header "TEST 9: Get Treatment Plan"
    
    if [ -z "$PLAN_VERSION_ID" ]; then
        print_error "No plan version ID available. Skipping."
        return
    fi
    
    response=$(curl -s -H "Authorization: Bearer $PATIENT_TOKEN" "$BASE_URL/api/v1/plans/$PLAN_VERSION_ID")
    
    if echo "$response" | grep -q "$PLAN_VERSION_ID"; then
        print_success "Treatment plan retrieved successfully"
        pretty_json "$response"
    else
        print_error "Failed to get treatment plan"
        echo "$response"
    fi
}

# ============================================
# TEST 10: Get Plan Estimate
# ============================================
test_get_estimate() {
    print_header "TEST 10: Get Plan Estimate"
    
    if [ -z "$PLAN_VERSION_ID" ]; then
        print_error "No plan version ID available. Skipping."
        return
    fi
    
    response=$(curl -s -H "Authorization: Bearer $PATIENT_TOKEN" "$BASE_URL/api/v1/plans/$PLAN_VERSION_ID/estimate")
    
    if echo "$response" | grep -q "estimates"; then
        print_success "Plan estimate retrieved"
        pretty_json "$response"
    else
        print_error "Failed to get estimate"
        echo "$response"
    fi
}

# ============================================
# TEST 11: Register Clinic
# ============================================
test_register_clinic() {
    print_header "TEST 11: Register Clinic"
    
    data='{
        "email": "clinic@test.com",
        "password": "password123",
        "role": "clinic_manager"
    }'
    
    response=$(api_call "POST" "/api/v1/auth/register" "$data")
    
    if echo "$response" | grep -q "access_token"; then
        print_success "Clinic user registered successfully"
        
        if [ "$JQ_AVAILABLE" = true ]; then
            CLINIC_TOKEN=$(echo "$response" | jq -r '.tokens.access_token')
        else
            CLINIC_TOKEN=$(echo "$response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        fi
        
        print_info "Clinic token: ${CLINIC_TOKEN:0:30}..."
    else
        print_error "Clinic registration failed"
        echo "$response"
    fi
}

# ============================================
# TEST 12: Create Clinic Profile
# ============================================
test_create_clinic_profile() {
    print_header "TEST 12: Create Clinic Profile"
    
    data='{
        "name": "Best Dental Clinic",
        "legal_name": "Best Dental LLC",
        "license_number": "LIC-12345",
        "year_established": 2015,
        "city": "Moscow",
        "district": "Central",
        "address": "Red Square, 1",
        "phone": "+7-495-123-4567",
        "email": "info@bestdental.com",
        "website": "https://bestdental.com",
        "price_segment": "business"
    }'
    
    response=$(api_call "POST" "/api/v1/clinic/profile" "$data" "$CLINIC_TOKEN")
    
    if echo "$response" | grep -q "license_number"; then
        print_success "Clinic profile created successfully"
        
        if [ "$JQ_AVAILABLE" = true ]; then
            CLINIC_ID=$(echo "$response" | jq -r '.id')
        else
            CLINIC_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | head -1)
        fi
        
        print_info "Clinic ID: $CLINIC_ID"
        pretty_json "$response"
    else
        print_error "Failed to create clinic profile"
        echo "$response"
    fi
}

# ============================================
# TEST 13: Add Items to Clinic Pricelist
# ============================================
test_add_pricelist_items() {
    print_header "TEST 13: Add Items to Clinic Pricelist"
    
    print_info "Adding therapy item..."
    data1='{
        "specialty": "therapy",
        "procedure_code": "D2391",
        "procedure_name": "Resin-based composite - one surface",
        "price_from": 5000,
        "price_to": 7000
    }'
    
    response1=$(api_call "POST" "/api/v1/clinic/pricelist" "$data1" "$CLINIC_TOKEN")
    
    if echo "$response1" | grep -q "specialty"; then
        print_success "Therapy item added"
    fi
    
    print_info "Adding orthopedics item..."
    data2='{
        "specialty": "orthopedics",
        "procedure_code": "D2740",
        "procedure_name": "Crown - porcelain/ceramic",
        "price_from": 25000,
        "price_to": 35000
    }'
    
    response2=$(api_call "POST" "/api/v1/clinic/pricelist" "$data2" "$CLINIC_TOKEN")
    
    if echo "$response2" | grep -q "specialty"; then
        print_success "Orthopedics item added"
    fi
    
    print_info "Adding surgery item..."
    data3='{
        "specialty": "surgery",
        "procedure_code": "D7140",
        "procedure_name": "Extraction, erupted tooth",
        "price_from": 3000,
        "price_to": 5000
    }'
    
    response3=$(api_call "POST" "/api/v1/clinic/pricelist" "$data3" "$CLINIC_TOKEN")
    
    if echo "$response3" | grep -q "specialty"; then
        print_success "Surgery item added"
        print_success "All pricelist items added successfully"
    else
        print_error "Failed to add some pricelist items"
    fi
}

# ============================================
# TEST 14: Get Clinic Pricelist
# ============================================
test_get_pricelist() {
    print_header "TEST 14: Get Clinic Pricelist"
    
    response=$(curl -s -H "Authorization: Bearer $CLINIC_TOKEN" "$BASE_URL/api/v1/clinic/pricelist")
    
    if echo "$response" | grep -q "specialty"; then
        print_success "Pricelist retrieved successfully"
        pretty_json "$response"
    else
        print_error "Failed to get pricelist"
        echo "$response"
    fi
}

# ============================================
# TEST 15: Create Offer Request
# ============================================
test_create_offer_request() {
    print_header "TEST 15: Create Offer Request (Patient)"
    
    if [ -z "$PLAN_VERSION_ID" ]; then
        print_error "No plan version ID available. Skipping."
        return
    fi
    
    # Get plan items
    plan_response=$(curl -s -H "Authorization: Bearer $PATIENT_TOKEN" "$BASE_URL/api/v1/plans/$PLAN_VERSION_ID")
    
    if [ "$JQ_AVAILABLE" = true ]; then
        item_ids=$(echo "$plan_response" | jq -r '[.plan_items[].id] | @json')
    else
        # Simplified: just use empty array if jq not available
        item_ids='[]'
        print_info "Using empty item selection (jq not available)"
    fi
    
    data='{
        "plan_version_id": "'$PLAN_VERSION_ID'",
        "selected_item_ids": '$item_ids',
        "preferred_city": "Moscow",
        "preferred_district": "Central",
        "price_segment": "business"
    }'
    
    response=$(api_call "POST" "/api/v1/patient/offer-requests" "$data" "$PATIENT_TOKEN")
    
    if echo "$response" | grep -q "plan_version_id"; then
        print_success "Offer request created successfully"
        
        if [ "$JQ_AVAILABLE" = true ]; then
            OFFER_REQUEST_ID=$(echo "$response" | jq -r '.id')
        else
            OFFER_REQUEST_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | head -1)
        fi
        
        print_info "Offer Request ID: $OFFER_REQUEST_ID"
        pretty_json "$response"
    else
        print_error "Failed to create offer request"
        echo "$response"
    fi
}

# ============================================
# MAIN EXECUTION
# ============================================

echo -e "${GREEN}"
echo "╔════════════════════════════════════════════╗"
echo "║  Dental Marketplace API Test Suite        ║"
echo "╚════════════════════════════════════════════╝"
echo -e "${NC}"

print_info "Testing against: $BASE_URL"

if [ "$JQ_AVAILABLE" = true ]; then
    print_success "jq is available - JSON output will be formatted"
else
    print_info "jq not found - JSON output will be raw (install with: sudo apt-get install jq)"
fi

echo ""

# Run tests
test_health
test_register_patient
test_login_patient
test_get_patient_profile
test_update_patient_profile
test_create_study
test_get_study
test_create_plan
test_get_plan
test_get_estimate
test_register_clinic
test_create_clinic_profile
test_add_pricelist_items
test_get_pricelist
test_create_offer_request

# Summary
print_header "TEST SUMMARY"
print_success "Tests completed!"
echo ""
print_info "Saved Variables:"
echo "  BASE_URL:      $BASE_URL"
echo "  Patient Token: ${PATIENT_TOKEN:0:30}..."
echo "  Clinic Token:  ${CLINIC_TOKEN:0:30}..."
echo "  Patient ID:    $PATIENT_ID"
echo "  Clinic ID:     $CLINIC_ID"
echo "  Study ID:      $STUDY_ID"
echo "  Plan ID:       $PLAN_VERSION_ID"
echo "  Offer Req ID:  $OFFER_REQUEST_ID"
echo ""

if [ -n "$CLINIC_ID" ]; then
    print_info "To activate clinic for testing offers:"
    echo "  docker exec dental_postgres psql -U postgres -d dental_marketplace -c \"UPDATE clinics SET is_active = true WHERE id = '$CLINIC_ID';\""
fi
