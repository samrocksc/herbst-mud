# World System

> 🔵 Last Updated: 2026-05-14  
> **Status**: 🟢 Complete

---

## Overview

The Herbst MUD engine supports **multi-world architecture** - the ability to run multiple independent MUD instances (worlds) from a single codebase. Each world has its own:
- Rooms, NPCs, Items, Quests, Abilities
- Dialog trees and content
- Player data (characters belong to users, not worlds)

---

## Architecture

### Database Schema

All world-scoped entities have a `world_id` field:

| Entity | Table Column | Description |
|--------|--------------|-------------|
| Room | `world_id` | Room belongs to a specific world |
| Ability | `world_id` | Ability is world-specific |
| NPCTemplate | `world_id` | NPC templates per world |
| EquipmentTemplate | `world_id` | Equipment templates per world |
| Quest | `world_id` | Quest definitions per world |
| DialogNode | `world_id` | Dialog trees per world |
| User | `allowed_worlds` | Comma-separated list of accessible worlds |

**Note**: Characters/Players do NOT have `world_id` - they belong to users, not worlds. Characters can access any world their user has permission for.

### Content Files

World content is organized in `content/worlds/`:

```
content/
└── worlds/
    ├── default/
    │   ├── rooms.yaml
    │   ├── npcs.yaml
    │   ├── items.yaml
    │   ├── quests.yaml
    │   └── dialog_trees.yaml
    ├── cyberpunk/
    │   └── ...
    └── steampunk/
        └── ...
```

---

## Multi-Tenancy: Whitelist Approach

World access is controlled via a **whitelist system**:

### User Whitelist

Each user has an `allowed_worlds` field containing a comma-separated list of world IDs:

```sql
-- Admin user (empty = all worlds)
allowed_worlds = ''

-- Regular user with limited access
allowed_worlds = 'cyberpunk,steampunk'
```

### Access Control

| User Type | `allowed_worlds` | Access |
|-----------|------------------|--------|
| Admin | Empty string | All worlds |
| User | `world1,world2` | Only specified worlds |
| User | Not set | All worlds (backward compatible) |

### Middleware Flow

```
Request
  ↓
AuthMiddleware (validate JWT)
  ↓
Fetch user.allowed_worlds from DB
  ↓
Store in context
  ↓
WorldAccessMiddleware
  ↓
Check if world_id in request matches whitelist
  ↓
Allow or 403 Forbidden
```

### Example API Request

```bash
# Access default world (admin)
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/rooms?world_id=default

# Access cyberpunk world
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/rooms?world_id=cyberpunk

# Forbidden if not in user's whitelist
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/rooms?world_id=restricted-world
```

---

## Configuration

### Creating a New World

1. **Add world to `content/worlds/`**:
   ```bash
   mkdir -p content/worlds/newworld
   # Copy existing content and customize
   ```

2. **Update database schema** (if needed):
   ```go
   // server/db/schema/your_entity.go
   field.String("world_id").Default("default")
   ```

3. **Run migrations**:
   ```bash
   cd server && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   ```

4. **Create users with access**:
   ```bash
   # Create user with access to specific worlds
   POST /api/users
   {
     "email": "user@example.com",
     "password": "password",
     "allowed_worlds": "newworld,cyberpunk"
   }
   ```

### Admin Configuration

Update user's allowed worlds:

```bash
PUT /api/users/:id
{
  "allowed_worlds": "cyberpunk,steampunk"
}
```

---

## Routes with World Access

All admin routes use `WorldAccessMiddleware`:

| Route | World Filtered |
|-------|----------------|
| `/api/rooms` | ✅ |
| `/api/abilities` | ✅ |
| `/api/npc-templates` | ✅ |
| `/api/equipment-templates` | ✅ |
| `/api/quests` | ✅ |
| `/api/dialog-nodes` | ✅ |

---

## Migration from Mono-World

To migrate from a mono-world setup:

1. **Set default world_id on existing data**:
   ```sql
   UPDATE rooms SET world_id = 'default' WHERE world_id IS NULL;
   UPDATE abilities SET world_id = 'default' WHERE world_id IS NULL;
   -- Repeat for other entities
   ```

2. **Create admin user with empty allowed_worlds**:
   ```sql
   UPDATE users SET allowed_worlds = NULL WHERE is_admin = true;
   ```

3. **Verify data**:
   ```bash
   # All rooms should have world_id
   SELECT COUNT(*) FROM rooms WHERE world_id IS NULL;
   # Should return 0
   ```

---

## Performance Notes

- `world_id` is indexed on all world-scoped tables
- Query filter: `WHERE world_id = ? OR world_id = 'default'`
- Empty `allowed_worlds` bypasses all checks (admin fast path)
- Context caching in middleware avoids DB hits per request

---

## Security

- Admin users bypass whitelist checks
- Non-admin users can only access whitelisted worlds
- Attempted access to non-whitelisted worlds returns `403 Forbidden`
- World ID validation happens before database queries

---

## Development

### Testing World Access

```bash
# Create test user with limited access
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test","allowed_worlds":"default"}'

# Test access to default world (should work)
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/rooms?world_id=default

# Test access to non-whitelisted world (should fail)
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/rooms?world_id=cyberpunk
```

---

## Related Docs

- [Effects System](./effects-system.md)
- [Character System Deep-dive](./CHARACTER-SYSTEM-DEEPDIVE.md)

---

🔵 Document version: 2026-05-14
