# Manual QA — Admin Panel Form Testing

## Quick Start
```bash
# 1. Source the test config
source _manual_qa_/config.sh

# 2. Login and get a token
TOKEN=$(curl -s -X POST "$API_URL/users/auth" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")

# 3. Test a specific feature
./_manual_qa_/test-npcs.sh
```

## Test Categories

| Feature | Script | Tests | Status |
|---------|--------|-------|--------|
| NPC Templates | `test-npcs.sh` | Create, list, edit, delete | ⏳ |
| Items | `test-items.sh` | Create, list, edit, delete | ⏳ |
| Abilities | `test-abilities.sh` | Create (with effects), list, edit, delete | ⏳ |
| Effects | `test-effects.sh` | Create, list, edit, delete, link to ability | ⏳ |
| Skills | `test-skills.sh` | Create, list, edit, delete | ⏳ |
| Quests | `test-quests.sh` | Create (with objectives), list, edit, delete | ⏳ |
| Triggers | `test-triggers.sh` | Create, list, edit, delete | ⏳ |
| Races | `test-races.sh` | Create, list, edit, delete | ⏳ |
| Genders | `test-genders.sh` | Create, list, edit, delete | ⏳ |
| Factions | `test-factions.sh` | Create, list, edit, delete | ⏳ |
| Players | `test-players.sh` | Create user, edit, reset password | ⏳ |
| Socials | `test-socials.sh` | Create, list | ⏳ |
| Channels | `test-channels.sh` | Create, list | ⏳ |
| Config | `test-config.sh` | Create, list, edit ($key), delete | ⏳ |
| Worlds | `test-worlds.sh` | Create, list, toggle active | ⏳ |
| Tags | `test-tags.sh` | Create, list | ⏳ |
| Map | `test-map.sh` | Create room, add exits, delete | ⏳ |
| Web Client | `test-web-client.sh` | Login, create character, room view, combat | ⏳ |

## Full Suite
```bash
./_manual_qa_/run-all.sh
```
