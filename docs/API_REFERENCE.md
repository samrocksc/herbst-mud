# API Reference

> 🔵 Last Updated: 2026-04-04

Complete REST API documentation for the Herbst MUD server.

**Base URL:**
- Development: `http://localhost:8080`
- Production: `https://your-api-domain.ondigitalocean.app`

**OpenAPI Spec:** Available at `/openapi.json`

---

## Table of Contents

- [Authentication](#authentication)
- [Health & Info](#health--info)
- [Users](#users)
- [Characters](#characters)
- [Rooms](#rooms)
- [Equipment](#equipment)
- [Skills & Talents](#skills--talents)
- [Quests](#quests)
- [Backups](#backups)
- [Error Handling](#error-handling)

---

## Authentication

Herbst MUD uses **JWT (JSON Web Tokens)** for authentication.

### Login

```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user_id": 1,
  "email": "user@example.com"
}
```

### Using the Token

Include the token in the `Authorization` header:

```http
GET /protected/resource
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### Token Validity

- **Expiration:** 24 hours from issuance
- **Algorithm:** HS256
- **Secret:** Configured via `JWT_SECRET` environment variable

---

## Health & Info

### Health Check

Check if the API is running.

```http
GET /healthz
```

**Response:**
```json
{
  "status": "ok",
  "ssh": "running",
  "db": "connected"
}
```

**Status Codes:**
- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service is down

### OpenAPI Specification

Returns the OpenAPI 3.0 spec for the API.

```http
GET /openapi.json
```

**Response:** `application/json` - Complete OpenAPI specification

---

## Users

### List All Users

```http
GET /users
```

**Authentication:** Required (Admin)

**Response:**
```json
[
  {
    "id": 1,
    "email": "admin@example.com",
    "is_admin": true
  },
  {
    "id": 2,
    "email": "player@example.com",
    "is_admin": false
  }
]
```

### Get User by ID

```http
GET /users/{id}
```

**Authentication:** Required

**Path Parameters:**
- `id` (integer) - User ID

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "is_admin": false
}
```

**Status Codes:**
- `200 OK` - Success
- `404 Not Found` - User doesn't exist

### Create User

```http
POST /users
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "securepassword",
  "isAdmin": false
}
```

**Authentication:** Not required (public registration)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | ✅ | User's email address |
| `password` | string | ✅ | Plain text password (hashed server-side) |
| `isAdmin` | boolean | ❌ | Admin privileges (default: false) |

**Response:**
```json
{
  "id": 3,
  "email": "newuser@example.com",
  "is_admin": false
}
```

**Status Codes:**
- `201 Created` - User created
- `400 Bad Request` - Invalid input
- `409 Conflict` - Email already exists

### Update User

```http
PUT /users/{id}
Content-Type: application/json

{
  "email": "updated@example.com",
  "password": "newpassword",
  "isAdmin": false
}
```

**Authentication:** Required

**Path Parameters:**
- `id` (integer) - User ID

**Request Body:** (All fields optional except password requirement)
- Same as Create User

**Response:** Updated user object

**Status Codes:**
- `200 OK` - Success
- `404 Not Found` - User doesn't exist

### Delete User

```http
DELETE /users/{id}
```

**Authentication:** Required (Admin or own account)

**Path Parameters:**
- `id` (integer) - User ID

**Status Codes:**
- `204 No Content` - Deleted successfully
- `404 Not Found` - User doesn't exist

---

## Characters

### List All Characters

```http
GET /characters
```

**Authentication:** Not required

**Response:**
```json
[
  {
    "id": 1,
    "name": "Gandalf",
    "level": 50,
    "hp": 100,
    "max_hp": 100,
    "room_id": 1,
    "user_id": 1
  }
]
```

### Get Character by ID

```http
GET /characters/{id}
```

**Authentication:** Not required

**Path Parameters:**
- `id` (integer) - Character ID

**Response:**
```json
{
  "id": 1,
  "name": "Gandalf",
  "level": 50,
  "hp": 100,
  "max_hp": 100,
  "strength": 15,
  "dexterity": 12,
  "constitution": 14,
  "intelligence": 18,
  "wisdom": 16,
  "charisma": 15,
  "room_id": 1,
  "user_id": 1,
  "skills": [...],
  "talents": [...],
  "inventory": [...],
  "equipment": {...}
}
```

### Create Character

```http
POST /characters
Content-Type: application/json

{
  "name": "Legolas",
  "level": 1,
  "strength": 14,
  "dexterity": 16,
  "constitution": 12,
  "intelligence": 10,
  "wisdom": 14,
  "charisma": 12,
  "user_id": 1
}
```

**Authentication:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Character name |
| `level` | integer | ❌ | Starting level (default: 1) |
| `strength` | integer | ❌ | STR stat |
| `dexterity` | integer | ❌ | DEX stat |
| `constitution` | integer | ❌ | CON stat |
| `intelligence` | integer | ❌ | INT stat |
| `wisdom` | integer | ❌ | WIS stat |
| `charisma` | integer | ❌ | CHA stat |
| `user_id` | integer | ✅ | Owner user ID |

**Response:**
```json
{
  "id": 2,
  "name": "Legolas",
  "level": 1,
  "hp": 100,
  "max_hp": 100,
  "room_id": 1
}
```

**Status Codes:**
- `201 Created` - Character created
- `400 Bad Request` - Invalid input

### Update Character

```http
PUT /characters/{id}
Content-Type: application/json

{
  "name": "Legolas the Wise",
  "level": 5
}
```

**Authentication:** Required (Owner)

**Path Parameters:**
- `id` (integer) - Character ID

**Status Codes:**
- `200 OK` - Success
- `403 Forbidden` - Not character owner
- `404 Not Found` - Character doesn't exist

### Delete Character

```http
DELETE /characters/{id}
```

**Authentication:** Required (Owner or Admin)

**Path Parameters:**
- `id` (integer) - Character ID

**Status Codes:**
- `204 No Content` - Deleted successfully
- `404 Not Found` - Character doesn't exist

### Character Login

```http
POST /characters/{id}/login
Authorization: Bearer {jwt_token}
```

**Authentication:** Required

**Description:** Authenticates a character for game session.

**Response:**
```json
{
  "character_id": 1,
  "token": "session_token...",
  "room": {
    "id": 1,
    "name": "Town Square",
    "description": "..."
  }
}
```

---

## Rooms

### List All Rooms

```http
GET /rooms
```

**Authentication:** Not required

**Response:**
```json
[
  {
    "id": 1,
    "name": "Town Square",
    "description": "A bustling town square...",
    "isStartingRoom": true,
    "exits": {
      "north": 2,
      "east": 3,
      "west": 4
    }
  }
]
```

### Get Room by ID

```http
GET /rooms/{id}
```

**Path Parameters:**
- `id` (integer) - Room ID

**Response:**
```json
{
  "id": 1,
  "name": "Town Square",
  "description": "A bustling town square...",
  "isStartingRoom": true,
  "exits": {
    "north": 2,
    "east": 3,
    "west": 4
  },
  "items": [...],
  "characters": [...]
}
```

### Create Room

```http
POST /rooms
Content-Type: application/json

{
  "name": "The Tavern",
  "description": "A cozy tavern...",
  "isStartingRoom": false,
  "exits": {
    "south": 1
  }
}
```

**Authentication:** Required (Admin)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Room name |
| `description` | string | ✅ | Room description |
| `isStartingRoom` | boolean | ❌ | New characters spawn here |
| `exits` | object | ❌ | Direction -> Room ID mapping |

### Update Room

```http
PUT /rooms/{id}
Content-Type: application/json

{
  "name": "Updated Tavern",
  "exits": {
    "south": 1,
    "up": 5
  }
}
```

**Authentication:** Required (Admin)

### Delete Room

```http
DELETE /rooms/{id}
```

**Authentication:** Required (Admin)

**Status Codes:**
- `204 No Content` - Deleted successfully
- `404 Not Found` - Room doesn't exist

---

## Equipment

### List Character Equipment

```http
GET /characters/{id}/equipment
```

**Authentication:** Required (Owner)

**Response:**
```json
{
  "head": null,
  "body": {
    "id": 1,
    "name": "Leather Armor",
    "type": "armor",
    "armor_class": 11
  },
  "hands": {
    "id": 2,
    "name": "Iron Sword",
    "type": "weapon",
    "damage": "1d8"
  },
  "feet": null
}
```

### Equip Item

```http
POST /characters/{id}/equip
Content-Type: application/json

{
  "item_id": 2,
  "slot": "hands"
}
```

**Authentication:** Required (Owner)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `item_id` | integer | ✅ | Item ID from inventory |
| `slot` | string | ✅ | Equipment slot (head, body, hands, feet, etc.) |

### Unequip Item

```http
POST /characters/{id}/unequip
Content-Type: application/json

{
  "slot": "hands"
}
```

**Authentication:** Required (Owner)

### Get Character Inventory

```http
GET /characters/{id}/inventory
```

**Authentication:** Required (Owner)

**Response:**
```json
[
  {
    "id": 3,
    "name": "Health Potion",
    "type": "consumable",
    "quantity": 3,
    "effects": [
      {
        "type": "heal",
        "amount": 25
      }
    ]
  }
]
```

### Use Consumable

```http
POST /characters/{id}/use-item
Content-Type: application/json

{
  "item_id": 3
}
```

**Authentication:** Required (Owner)

---

## Skills & Talents

### Get Character Skills

```http
GET /characters/{id}/skills
```

**Authentication:** Required (Owner)

**Response:**
```json
[
  {
    "id": 1,
    "name": "Sword Mastery",
    "level": 3,
    "max_level": 5,
    "description": "Increases sword damage"
  }
]
```

### Get Available Talents

```http
GET /characters/{id}/available-talents
```

**Authentication:** Required (Owner)

**Response:**
```json
[
  {
    "id": 1,
    "name": "Power Attack",
    "description": "Next attack deals double damage",
    "cooldown": 3,
    "unlocked": true
  }
]
```

### Equip Talent

```http
POST /characters/{id}/talents/equip
Content-Type: application/json

{
  "talent_id": 1,
  "slot": 1
}
```

**Authentication:** Required (Owner)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `talent_id` | integer | ✅ | Talent ID |
| `slot` | integer | ✅ | Talent slot (1-6) |

### Unequip Talent

```http
POST /characters/{id}/talents/unequip
Content-Type: application/json

{
  "slot": 1
}
```

**Authentication:** Required (Owner)

---

## Quests

### Get Active Quests

```http
GET /characters/{id}/quests
```

**Authentication:** Required (Owner)

**Response:**
```json
[
  {
    "id": 1,
    "name": "The Beginning",
    "description": "Find the Fountain of Dreams",
    "status": "in_progress",
    "objectives": [
      {
        "id": 1,
        "description": "Visit the Fountain",
        "completed": false
      }
    ],
    "rewards": {
      "xp": 100,
      "gold": 50
    }
  }
]
```

### Get Quest by ID

```http
GET /quests/{id}
```

**Authentication:** Required

### List All Quests

```http
GET /quests
```

**Authentication:** Required (Admin)

---

## Backups

### Export Database

```http
GET /backup/export
```

**Authentication:** Required (Admin)

**Response:** `application/json` - Full database export

**Description:** Exports all rooms, items, quests, and configuration as JSON.

### Import Database

```http
POST /backup/import
Content-Type: application/json

{
  "rooms": [...],
  "items": [...],
  "quests": [...]
}
```

**Authentication:** Required (Admin)

**Warning:** Import can overwrite existing data.

---

## Error Handling

### Error Response Format

All errors follow this format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional context"
  }
}
```

### HTTP Status Codes

| Code | Meaning | Common Causes |
|------|---------|---------------|
| `200 OK` | Success | Request succeeded |
| `201 Created` | Created | Resource created successfully |
| `204 No Content` | No Content | Delete succeeded |
| `400 Bad Request` | Bad Request | Invalid JSON, missing fields |
| `401 Unauthorized` | Unauthorized | Missing/invalid JWT token |
| `403 Forbidden` | Forbidden | Insufficient permissions |
| `404 Not Found` | Not Found | Resource doesn't exist |
| `409 Conflict` | Conflict | Duplicate email, etc. |
| `429 Too Many Requests` | Rate Limited | Too many requests |
| `500 Internal Server Error` | Server Error | Unexpected error |
| `503 Service Unavailable` | Unavailable | Service down |

### Rate Limiting Headers

When rate limiting is triggered:

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1640000000
Retry-After: 60
```

---

## Code Examples

### cURL Examples

**Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

**Create Character:**
```bash
curl -X POST http://localhost:8080/characters \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"Frodo","user_id":1}'
```

**List Rooms:**
```bash
curl http://localhost:8080/rooms
```

### JavaScript/TypeScript Example

```typescript
const BASE_URL = 'http://localhost:8080';

// Login
async function login(email: string, password: string) {
  const response = await fetch(`${BASE_URL}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  return response.json();
}

// Get character with auth
async function getCharacter(id: number, token: string) {
  const response = await fetch(`${BASE_URL}/characters/${id}`, {
    headers: { 'Authorization': `Bearer ${token}` },
  });
  return response.json();
}
```

### Go Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func login(email, password string) (*LoginResponse, error) {
    payload := map[string]string{
        "email":    email,
        "password": password,
    }
    data, _ := json.Marshal(payload)
    
    resp, err := http.Post(
        "http://localhost:8080/login",
        "application/json",
        bytes.NewBuffer(data),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result LoginResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

---

## WebSocket/Real-time

Currently, the game uses HTTP polling for real-time updates. WebSocket support is planned for a future release.

---

🔵 Document version: 2026-04-04
