#!/bin/bash

# QSO API Test Script for Ham-Radio Cloud

set -e

API_URL="${API_URL:-http://localhost:8080}"
BASE_URL="$API_URL/api/v1"

echo "🧪 Testing Ham-Radio Cloud QSO API"
echo "   API URL: $API_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to print test results
test_endpoint() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expected_status="$5"

    echo -e "${BLUE}Testing:${NC} $name"
    echo "   $method $endpoint"

    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "   ${GREEN}✓ PASS${NC} (HTTP $http_code)"
        ((TESTS_PASSED++))
        echo "$body" | jq . 2>/dev/null || echo "$body"
    else
        echo -e "   ${RED}✗ FAIL${NC} (Expected HTTP $expected_status, got $http_code)"
        ((TESTS_FAILED++))
        echo "$body"
    fi
    echo ""
}

# 1. Health Check
test_endpoint "Health Check" "GET" "/health" "" 200

# 2. API Info
test_endpoint "API Info" "GET" "/" "" 200

# 3. Get QSO Stats (should work even with no QSOs)
test_endpoint "Get QSO Stats" "GET" "/qsos/stats" "" 200

# 4. List QSOs (empty initially)
test_endpoint "List QSOs (empty)" "GET" "/qsos" "" 200

# 5. Create a QSO
QSO_DATA='{
    "callsign": "K1ABC",
    "time_on": "2025-11-23T14:30:00Z",
    "band": "20m",
    "mode": "SSB",
    "rst_sent": "59",
    "rst_rcvd": "57",
    "freq": 14.250,
    "name": "John",
    "qth": "Boston, MA",
    "gridsquare": "FN42",
    "tx_pwr": 100,
    "comment": "Nice contact!"
}'

response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d "$QSO_DATA" \
    "$BASE_URL/qsos")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo -e "${BLUE}Testing:${NC} Create QSO"
echo "   POST /qsos"

if [ "$http_code" -eq 201 ]; then
    echo -e "   ${GREEN}✓ PASS${NC} (HTTP $http_code)"
    ((TESTS_PASSED++))
    echo "$body" | jq .
    QSO_ID=$(echo "$body" | jq -r '.data.id')
    echo "   Created QSO ID: $QSO_ID"
else
    echo -e "   ${RED}✗ FAIL${NC} (Expected HTTP 201, got $http_code)"
    ((TESTS_FAILED++))
    echo "$body"
    QSO_ID=""
fi
echo ""

# 6. Get the created QSO
if [ -n "$QSO_ID" ]; then
    test_endpoint "Get QSO by ID" "GET" "/qsos/$QSO_ID" "" 200
fi

# 7. Update the QSO
if [ -n "$QSO_ID" ]; then
    UPDATE_DATA='{
        "callsign": "K1ABC",
        "time_on": "2025-11-23T14:30:00Z",
        "time_off": "2025-11-23T14:45:00Z",
        "band": "20m",
        "mode": "SSB",
        "rst_sent": "59",
        "rst_rcvd": "59",
        "freq": 14.250,
        "name": "John Smith",
        "qth": "Boston, MA",
        "gridsquare": "FN42",
        "tx_pwr": 100,
        "comment": "Updated: Great QSO!"
    }'
    test_endpoint "Update QSO" "PUT" "/qsos/$QSO_ID" "$UPDATE_DATA" 200
fi

# 8. List QSOs (should have 1 now)
test_endpoint "List QSOs (with data)" "GET" "/qsos" "" 200

# 9. Filter QSOs by band
test_endpoint "Filter QSOs by band" "GET" "/qsos?band=20m" "" 200

# 10. Filter QSOs by callsign
test_endpoint "Filter QSOs by callsign" "GET" "/qsos?callsign=K1ABC" "" 200

# 11. Get QSO Stats (should show 1 QSO)
test_endpoint "Get QSO Stats (updated)" "GET" "/qsos/stats" "" 200

# 12. Create another QSO
QSO_DATA_2='{
    "callsign": "N2DEF",
    "time_on": "2025-11-23T15:00:00Z",
    "band": "40m",
    "mode": "CW",
    "rst_sent": "599",
    "rst_rcvd": "579",
    "freq": 7.030,
    "gridsquare": "FN31",
    "tx_pwr": 50
}'
test_endpoint "Create 2nd QSO" "POST" "/qsos" "$QSO_DATA_2" 201

# 13. List QSOs with pagination
test_endpoint "List QSOs (limit=1)" "GET" "/qsos?limit=1" "" 200

# 14. Delete the first QSO
if [ -n "$QSO_ID" ]; then
    test_endpoint "Delete QSO" "DELETE" "/qsos/$QSO_ID" "" 204
fi

# 15. Try to get deleted QSO (should fail)
if [ -n "$QSO_ID" ]; then
    test_endpoint "Get deleted QSO (should fail)" "GET" "/qsos/$QSO_ID" "" 404
fi

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Test Summary:"
echo -e "   ${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "   ${RED}Failed: $TESTS_FAILED${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi
