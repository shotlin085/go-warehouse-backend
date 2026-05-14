#!/bin/bash

echo "🧪 Testing Go Warehouse API"
echo "=============================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
echo "GET http://localhost:7000/health"
curl -s http://localhost:7000/health | jq .
echo ""
echo ""

# Test 2: Create Item
echo -e "${YELLOW}Test 2: Create Item${NC}"
echo "POST http://localhost:7000/api/items"
RESPONSE=$(curl -s -X POST http://localhost:7000/api/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Widget",
    "sku": "TEST-001",
    "description": "A test item for school project",
    "quantity": 10,
    "unit_price": 25.99
  }')
echo "$RESPONSE" | jq .
ITEM_ID=$(echo "$RESPONSE" | jq -r '.id')
echo ""
echo ""

# Test 3: Get All Items
echo -e "${YELLOW}Test 3: Get All Items${NC}"
echo "GET http://localhost:7000/api/items"
curl -s http://localhost:7000/api/items | jq .
echo ""
echo ""

# Test 4: Get Single Item
if [ ! -z "$ITEM_ID" ] && [ "$ITEM_ID" != "null" ]; then
  echo -e "${YELLOW}Test 4: Get Single Item (ID: $ITEM_ID)${NC}"
  echo "GET http://localhost:7000/api/items/$ITEM_ID"
  curl -s http://localhost:7000/api/items/$ITEM_ID | jq .
  echo ""
  echo ""
fi

echo -e "${GREEN}✅ All tests completed!${NC}"
