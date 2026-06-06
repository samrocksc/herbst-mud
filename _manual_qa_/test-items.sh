#!/usr/bin/env bash
# Manual QA: Items CRUD
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
source "$DIR/config.sh"

echo "=== Items QA ==="
check_server
TOKEN=*** + PASSWORD="$PASSWORD"; PASS=0; FAIL=0
TEST_NAME="QA-Item-$(date +%s)"

echo -e "\n${YELLOW}1. CREATE${NC}"
RESP=$(api POST "$API_URL/api/equipment-templates?world_id=default" \
  "{\"name\":\"$TEST_NAME\",\"description\":\"Test sword\",\"slot\":\"main_hand\",\"item_type\":\"weapon\",\"level\":1,\"weight\":5,\"rarity\":\"common\",\"armor_rating\":0,\"armor_type\":\"\",\"stats\":{},\"damage_dice_count\":1,\"damage_dice_sides\":8,\"damage_bonus\":1,\"damage_type\":\"slashing\",\"weapon_type\":\"sword\",\"is_two_handed\":false,\"is_visible\":true,\"is_immovable\":false,\"is_container\":false}")
HTTP=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | head -1)
ITEM_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('id','') or d.get('equipment_template',{}).get('id',''))" 2>/dev/null)

if assert_status "201" "$HTTP" "Create Item"; then PASS=$((PASS+1))
else FAIL=$((FAIL+1)); echo "Body: $(echo "$BODY" | head -c 300)"; fi

if [ -z "$ITEM_ID" ]; then echo -e "${RED}FAIL${NC} No ID"; echo "Results: $PASS/$FAIL"; exit 1; fi
echo "Created: $ITEM_ID"

echo -e "\n${YELLOW}2. LIST${NC}"
RESP=$(api GET "$API_URL/api/equipment-templates?world_id=default")
HTTP=$(echo "$RESP" | tail -1)
if assert_status "200" "$HTTP" "List Items"; then PASS=$((PASS+1)); else FAIL=$((FAIL+1)); fi

echo -e "\n${YELLOW}3. EDIT${NC}"
RESP=$(api PUT "$API_URL/api/equipment-templates/$ITEM_ID" \
  "{\"name\":\"$TEST_NAME-edited\",\"description\":\"edited\",\"slot\":\"main_hand\",\"item_type\":\"weapon\",\"level\":5,\"weight\":8,\"rarity\":\"rare\",\"armor_rating\":0,\"armor_type\":\"\",\"stats\":{},\"damage_dice_count\":2,\"damage_dice_sides\":6,\"damage_bonus\":3,\"damage_type\":\"fire\",\"weapon_type\":\"sword\",\"is_two_handed\":false,\"is_visible\":true,\"is_immovable\":false,\"is_container\":false}")
HTTP=$(echo "$RESP" | tail -1)
if assert_status "200" "$HTTP" "Edit Item"; then PASS=$((PASS+1)); else FAIL=$((FAIL+1)); fi

echo -e "\n${YELLOW}4. DELETE${NC}"
RESP=$(api DELETE "$API_URL/api/equipment-templates/$ITEM_ID")
HTTP=$(echo "$RESP" | tail -1)
if assert_status "200" "$HTTP" "Delete Item"; then PASS=$((PASS+1))
  RESP=$(api GET "$API_URL/api/equipment-templates?world_id=default")
  BODY=$(echo "$RESP" | head -1)
  FOUND=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); items=d if isinstance(d,list) else d.get('equipment_templates',[]); print(any(str(i.get('id'))=='$ITEM_ID' for i in items))")
  if [ "$FOUND" != "True" ]; then echo -e "  ${GREEN}PASS${NC} Gone from list"; PASS=$((PASS+1))
  else echo -e "  ${RED}FAIL${NC} Still in list"; FAIL=$((FAIL+1)); fi
else FAIL=$((FAIL+1)); fi

echo -e "\n${YELLOW}=== Items Results: $PASS passed, $FAIL failed ===${NC}"
[ "$FAIL" -eq 0 ]
