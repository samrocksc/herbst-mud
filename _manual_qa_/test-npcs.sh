#!/usr/bin/env bash
# Manual QA: NPC Template CRUD
# Usage: ./_manual_qa_/test-npcs.sh

set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
source "$DIR/config.sh"

echo "=== NPC Template QA ==="
check_server

TOKEN=***  
TEST_NAME="QA-NPC-$(date +%s)"
PASS=0
FAIL=0

# --- CREATE ---
echo -e "\n${YELLOW}1. CREATE${NC}"
RESP=$(api POST "$API_URL/api/npc-templates?world_id=default" \
  "{\"name\":\"$TEST_NAME\",\"description\":\"Manual QA test\",\"race_id\":0,\"disposition\":\"neutral\",\"level\":1,\"xp_value\":50,\"respawn_cooldown\":60,\"respawn_rooms\":[],\"greeting\":\"Hello\",\"skills\":{},\"trades_with\":[]}")

HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | head -1)
NPC_ID=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))" 2>/dev/null)

if assert_status "201" "$HTTP" "Create NPC Template ($TEST_NAME)"; then
    PASS=$((PASS+1))
else
    FAIL=$((FAIL+1))
    echo "  Response: $(echo "$BODY" | head -c 200)"
fi

if [ -z "$NPC_ID" ] || [ "$NPC_ID" = "null" ]; then
    echo -e "${RED}FAIL${NC} No ID returned — skipping remaining tests"
    echo "Results: $PASS passed, $FAIL failed"
    exit 1
fi
echo "  Created: $NPC_ID"

# --- LIST / Verify ---
echo -e "\n${YELLOW}2. LIST${NC}"
RESP=$(api GET "$API_URL/api/npc-templates?world_id=default")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | head -1)

if assert_status "200" "$HTTP" "List NPC Templates"; then
    FOUND=$(echo "$BODY" | python3 -c "import sys,json; items=json.load(sys.stdin); print(any(i.get('name')=='$TEST_NAME' for i in (items if isinstance(items,list) else items.get('npcs',[]))))")
    if [ "$FOUND" = "True" ]; then
        echo -e "  ${GREEN}PASS${NC} Created NPC found in list"
        PASS=$((PASS+1))
    else
        echo -e "  ${YELLOW}WARN${NC} Created NPC not in list response"
    fi
else
    FAIL=$((FAIL+1))
fi

# --- EDIT ---
echo -e "\n${YELLOW}3. EDIT${NC}"
RESP=$(api PUT "$API_URL/api/npc-templates/$NPC_ID" \
  "{\"name\":\"$TEST_NAME-edited\",\"description\":\"Edited\",\"race_id\":0,\"disposition\":\"hostile\",\"level\":5,\"xp_value\":100,\"respawn_cooldown\":120,\"respawn_rooms\":[],\"greeting\":\"Grrr\",\"skills\":{},\"trades_with\":[]}")
HTTP=$(echo "$RESP" | tail -1)

if assert_status "200" "$HTTP" "Edit NPC Template"; then
    PASS=$((PASS+1))
else
    FAIL=$((FAIL+1))
fi

# --- DELETE ---
echo -e "\n${YELLOW}4. DELETE${NC}"
RESP=$(api DELETE "$API_URL/api/npc-templates/$NPC_ID")
HTTP=$(echo "$RESP" | tail -1)

if assert_status "200" "$HTTP" "Delete NPC Template"; then
    PASS=$((PASS+1))
    
    # Verify gone
    RESP=$(api GET "$API_URL/api/npc-templates?world_id=default")
    BODY=$(echo "$RESP" | head -1)
    FOUND=$(echo "$BODY" | python3 -c "import sys,json; items=json.load(sys.stdin); print(any(i.get('id')=='$NPC_ID' for i in (items if isinstance(items,list) else items.get('npcs',[]))))")
    if [ "$FOUND" = "False" ]; then
        echo -e "  ${GREEN}PASS${NC} Deleted NPC no longer in list"
        PASS=$((PASS+1))
    else
        echo -e "  ${RED}FAIL${NC} Deleted NPC still appears in list"
        FAIL=$((FAIL+1))
    fi
else
    FAIL=$((FAIL+1))
fi

# --- Edge Cases ---
echo -e "\n${YELLOW}5. EDGE CASES${NC}"

# Missing name
RESP=$(api POST "$API_URL/api/npc-templates?world_id=default" \
  '{"name":"","description":"","race_id":0,"disposition":"neutral","level":1,"xp_value":0,"respawn_cooldown":60,"respawn_rooms":[],"greeting":"","skills":{},"trades_with":[]}')
HTTP=$(echo "$RESP" | tail -1)
if assert_status "400" "$HTTP" "Reject empty name"; then PASS=$((PASS+1)); else FAIL=$((FAIL+1)); fi

# Non-existent edit
RESP=$(api PUT "$API_URL/api/npc-templates/non-existent-id" '{"name":"ghost"}')
HTTP=$(echo "$RESP" | tail -1)
if assert_status "404" "$HTTP" "404 on non-existent edit"; then PASS=$((PASS+1)); else FAIL=$((FAIL+1)); fi

# --- Summary ---
echo -e "\n${YELLOW}=== NPC Template Results: $PASS passed, $FAIL failed ===${NC}"
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
