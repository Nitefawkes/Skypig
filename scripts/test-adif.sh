#!/bin/bash

# ADIF Import/Export Test Script for Ham-Radio Cloud

set -e

API_URL="${API_URL:-http://localhost:8080}"
BASE_URL="$API_URL/api/v1"

echo "🧪 Testing Ham-Radio Cloud ADIF Import/Export"
echo "   API URL: $API_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to print test results
test_endpoint() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local data_file="$4"
    local expected_status="$5"

    echo -e "${BLUE}Testing:${NC} $name"
    echo "   $method $endpoint"

    if [ -z "$data_file" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: text/plain" \
            --data-binary "@$data_file" \
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

# Check if sample ADIF file exists
SAMPLE_FILE="backend/testdata/sample.adi"
if [ ! -f "$SAMPLE_FILE" ]; then
    echo -e "${RED}✗ Sample ADIF file not found: $SAMPLE_FILE${NC}"
    exit 1
fi

echo -e "${BLUE}Using test file:${NC} $SAMPLE_FILE"
echo ""

# 1. Validate ADIF file
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Part 1: Validation${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

test_endpoint "Validate ADIF file" "POST" "/qsos/validate" "$SAMPLE_FILE" 200

# 2. Import ADIF file
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Part 2: Import${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${BLUE}Testing:${NC} Import ADIF file"
echo "   POST /qsos/import"

response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: text/plain" \
    --data-binary "@$SAMPLE_FILE" \
    "$BASE_URL/qsos/import")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 206 ]; then
    echo -e "   ${GREEN}✓ PASS${NC} (HTTP $http_code)"
    ((TESTS_PASSED++))
    echo "$body" | jq .

    IMPORTED=$(echo "$body" | jq -r '.data.imported_records')
    echo -e "${GREEN}   Imported $IMPORTED records${NC}"
else
    echo -e "   ${RED}✗ FAIL${NC} (Expected HTTP 200/206, got $http_code)"
    ((TESTS_FAILED++))
    echo "$body"
fi
echo ""

# 3. Verify QSOs were imported
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Part 3: Verify Import${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

test_endpoint "List imported QSOs" "GET" "/qsos?limit=10" "" 200

# 4. Export all QSOs
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Part 4: Export${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${BLUE}Testing:${NC} Export all QSOs to ADIF"
echo "   GET /qsos/export"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/qsos/export")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
    echo -e "   ${GREEN}✓ PASS${NC} (HTTP $http_code)"
    ((TESTS_PASSED++))

    # Count records in export
    record_count=$(echo "$body" | grep -c "<eor>" || echo "0")
    echo -e "${GREEN}   Exported $record_count records${NC}"

    # Show first few lines
    echo ""
    echo -e "${YELLOW}Export preview:${NC}"
    echo "$body" | head -20
    echo "   [... truncated ...]"
else
    echo -e "   ${RED}✗ FAIL${NC} (Expected HTTP 200, got $http_code)"
    ((TESTS_FAILED++))
fi
echo ""

# 5. Export with filter (band=20m)
echo -e "${BLUE}Testing:${NC} Export filtered QSOs (band=20m)"
echo "   GET /qsos/export?band=20m"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/qsos/export?band=20m")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
    echo -e "   ${GREEN}✓ PASS${NC} (HTTP $http_code)"
    ((TESTS_PASSED++))

    record_count=$(echo "$body" | grep -c "<eor>" || echo "0")
    echo -e "${GREEN}   Exported $record_count records on 20m${NC}"
else
    echo -e "   ${RED}✗ FAIL${NC} (Expected HTTP 200, got $http_code)"
    ((TESTS_FAILED++))
fi
echo ""

# 6. Test invalid ADIF
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Part 5: Error Handling${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${BLUE}Testing:${NC} Validate invalid ADIF"
echo "   POST /qsos/validate"

INVALID_ADIF="<CALL:5>W1AW <eor>"
response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: text/plain" \
    -d "$INVALID_ADIF" \
    "$BASE_URL/qsos/validate")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

# Should return 400 for invalid ADIF
if [ "$http_code" -eq 400 ]; then
    echo -e "   ${GREEN}✓ PASS${NC} (HTTP $http_code - correctly rejected invalid ADIF)"
    ((TESTS_PASSED++))
    echo "$body" | jq .
else
    echo -e "   ${YELLOW}⚠ WARNING${NC} (Expected HTTP 400, got $http_code)"
    echo "$body"
fi
echo ""

# 7. Test empty import
echo -e "${BLUE}Testing:${NC} Import empty content"
echo "   POST /qsos/import"

response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: text/plain" \
    -d "" \
    "$BASE_URL/qsos/import")

http_code=$(echo "$response" | tail -n1)

if [ "$http_code" -eq 400 ]; then
    echo -e "   ${GREEN}✓ PASS${NC} (HTTP $http_code - correctly rejected empty content)"
    ((TESTS_PASSED++))
else
    echo -e "   ${RED}✗ FAIL${NC} (Expected HTTP 400, got $http_code)"
    ((TESTS_FAILED++))
fi
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Test Summary:"
echo -e "   ${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "   ${RED}Failed: $TESTS_FAILED${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All ADIF tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi
