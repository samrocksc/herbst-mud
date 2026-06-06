#!/usr/bin/env bash
# Manual QA: Effects CRUD
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
source "$DIR/config.sh"

echo "=== Effects QA ==="
check_server

TOKEN=***  
TEST_NAME="QA-Effect-$(date +%s)"
PASS=0
FAIL=0

echo -e "\n${YELLOW}1. CREATE${NC}"
RESP=$(api POST "$API_URL/api/effects" \
  "{\"name\":\"$TEST_NAME\",\"description\":\"QA test damage effect\",\"effect_type\":\"hp_change\",\"parameters\":{\"amount\":-50},\"stack_mode\":\"replace\",\"stack_limit\":1,\"is_permanent\":false,\"duration_secs\":0,\"messages\":{\"on_start\":\"You burn!\",\"on_end\":\"Fire fades\"}}")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | head -1)
EFFECT_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('id','') or d.get('effect',{}).get('id',''))" 2>/dev/null)

if assert_status "201" "$HTTP" "Create Effect"; then PASS=$((PASS+1))
else FAIL=$((FAIL+1)); echo "  Body: $(echo "$BODY" | head -c 300)"; fi

if [ -z "$EFFECT_ID" ] || [ "$EFFECT_ID" = "null" ]; then
    echo -e "${RED}FAIL${NC} No effect ID returned"
    echo "Results: $PASS passed, $FAIL failed"
    exit 1
fi
echo "  Created effect ID: $EFFECT_ID"

echo -e "\n${YELLOW}2. LIST${NC}"
RESP=$(api GET "$API_URL/api/effects")
HTTP=$(echo "$RESP" | tail -1)
if assert_status "200" "$HTTP" "List Effects"; then PASS=$((PASS+1)); else FAIL=$((FAIL+1)); fi

echo -e "\n${YELLOW}3. EDIT${NC}"
RESP=$(api PUT "$API_URL/api/effects/$EFFECT_ID" \
  "{\"name\":\"$TEST_NAME-edited\",\"description\":\"Edited\",\"effect_type\":\"hp_change\",\"parameters\":{\"amount\":-99},\"stack_mode\":\"replace\",\"stack_limit\":1,\"is_permanent\":false,\"duration_secs\":0,\"messages\":{\"on_start\":\"OUCH\",\"on_end\":\"ok\"}}")
HTTP=$(echo "$RESP" | tail -1)
if assert_status "200" "$HTTP" "Edit Effect"; then PASS=$((PASS+1)); else FAIL=$((FAIL+1)); fi

echo -e "\n${YELLOW}4. DELETE${NC}"
RESP=$(api DELETE "$API_URL/api/effects/$EFFECT_ID")
HTTP=$(echo "$RESP" | tail -1)
if assert_status "200" "$HTTP" "Delete Effect"; then PASS=$((PASS+1)); else FAIL=$((FAIL+1)); fi

echo -e "\n${YELLOW}5. LINK TO ABILITY${NC}"
# Create effect, create ability, link them
EFF_RESP=$(api POST "$API_URL/api/effects" \
  "{\"name\":\"QA-Ability-Effect\",\"description\":\"For ability test\",\"effect_type\":\"hp_change\",\"parameters\":{\"amount\":-25},\"stack_mode\":\"replace\",\"stack_limit\":1,\"is_permanent\":false,\"duration_secs\":0,\"messages\":{}}")
EFF_HTTP=$(echo "$EFF_RESP" | tail -1)
EFF_BODY=$(echo "$EFF_RESP" | head -1)
EFF_ID=$(echo "$EFF_BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('id','') or d.get('effect',{}).get('id',''))" 2>/dev/null)

ABL_RESP=$(api POST "$API_URL/api/abilities?world_id=default" \
  "{\"name\":\"QA-Ability-Test\",\"description\":\"Test\",\"ability_type\":\"combat\",\"ability_class\":\"active\",\"requirements\":\"1\",\"cost\":0,\"cooldown_seconds\":0,\"mana_cost\":0,\"stamina_cost\":0,\"hp_cost\":0,\"proc_chance\":0,\"proc_event\":\"\"}")
ABL_HTTP=$(echo "$ABL_RESP" | tail -1)
ABL_BODY=$(echo "$ABL_RESP" | head -1)
ABL_ID=$(echo "$ABL_BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('id',''))" 2>/dev/null)

if [ -n "$ABL_ID" ] && [ -n "$EFF_ID" ]; then
    LINK_RESP=$(api POST "$API_URL/api/abilities/$ABL_ID/effects" \
      "{\"effect_id\":$EFF_ID,\"effect_type\":\"hp_change\",\"target\":\"enemy\",\"value\":-25,\"duration\":0,\"scaling_stat\":\"\",\"scaling_ratio\":0,\"sort_order\":0,\"effect_message\":\"Take 25 damage!\"}")
    LINK_HTTP=$(echo "$LINK_RESP" | tail -1)
    if assert_status "201" "$LINK_HTTP" "Link effect to ability"; then PASS=$((PASS+1))
    else FAIL=$((FAIL+1)); fi
    
    # Cleanup
    api DELETE "$API_URL/api/abilities/$ABL_ID/effects/$EFF_ID" > /dev/null 2>&1
    api DELETE "$API_URL/api/abilities/$ABL_ID" > /dev/null 2>&1
fi
api DELETE "$API_URL/api/effects/$EFF_ID" > /dev/null 2>&1

echo -e "\n${YELLOW}=== Effects Results: $PASS passed, $FAIL failed ===${NC}"
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
