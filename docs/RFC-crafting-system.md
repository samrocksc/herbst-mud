# RFC: Crafting System — Stations, Recipes, and Item Crafting

**Status:** Draft  
**Date:** 2026-05-19  
**Author:** Leonardo (Donnie)  
**Analog:** Pizza Chef class — craft a pizza, wield it as a weapon

---

## 1. Goal

Enable players to craft items at in-world stations. The primary analog is a **Pizza Chef** class:
- Stand at a `pizza_station` in a room
- Have the right ingredients in inventory (dough, sauce, cheese, toppings)
- Execute `craft pepperoni_pizza` → item appears in inventory
- Equip and wield the pizza as a weapon (or eat it, or throw it)

The system should be generic enough for any class/recipe: blacksmithing, alchemy, cooking, tinkering.

---

## 2. Design Philosophy

- **Engine, not content.** The crafting engine defines the schema and mechanics. Recipes, stations, and ingredient types are content created through the admin UI.
- **Reuse existing patterns.** Template → Instance, Effect system, Faction categories, Tag gating — all already exist.
- **No new item type.** "Ingredient" is a tag, not a new enum value. An item is a weapon when wielded regardless of what it's made of.
- **Room-tag based discovery.** Stations are identified by tags on the room, not separate entities.

---

## 3. New Schema Entities

### 3.1 Recipe

```go
// server/db/schema/crafting_recipe.go
type CraftingRecipe struct {
    ent.Schema
}

func (CraftingRecipe) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Unique(),          // "pepperoni_pizza"
        field.String("display_name"),             // "Pepperoni Pizza"
        field.Text("description").Optional(),
        field.String("required_station_tag"),     // "pizza_station"
        field.String("required_class").Optional(), // faction-based class name, e.g. "pizza_chef"
        field.Int("required_skill_level").Default(0), // minimum skill level (e.g. cooking)
        field.String("required_skill").Optional(),    // skill name e.g. "cooking"
        field.JSON("inputs", []CraftingInput{}),  // what goes in
        field.JSON("outputs", []CraftingOutput{}), // what comes out
        field.Int("craft_time_secs").Default(3),  // channel time
        field.String("world_id").Default("default"),
    }
}

type CraftingInput struct {
    EquipmentTemplateID string `json:"equipment_template_id"` // required
    Quantity            int    `json:"quantity"`               // default 1
    Consumed            bool   `json:"consumed"`               // true = destroyed on craft
}

type CraftingOutput struct {
    EquipmentTemplateID string `json:"equipment_template_id"`
    Quantity            int    `json:"quantity"` // default 1
}
```

**Why not an M2M edge to equipment_templates?** Inputs/outputs use template IDs stored as JSON because:
- A recipe can reference any number of templates
- JSON is simpler than join tables for this first iteration
- Can be extended with `tag_filter` later (e.g. "any item tagged `metal`")

### 3.2 Station (Room Tag approach — NO new entity)

Stations are **room tags** — a simple naming convention:

| Room Tag | Station Type |
|----------|-------------|
| `pizza_station` | Pizza oven / prep table |
| `forge` | Blacksmith forge |
| `alchemy_table` | Potion crafting |
| `workbench` | General tinkering |

A room can have multiple station tags. The `craft` command checks if the current room has the tag matching `required_station_tag` on the recipe.

**Implementation:** Add a `tags` JSON field to the Room schema, or add room edges to the existing Tag entity.

### 3.3 Item Quantity

Add a `quantity` field to Equipment:

```go
field.Int("quantity").Default(1).Comment("Stack size for consumable/ingredient items")
```

This avoids 2000 separate rows for 2000 units of dough. When quantity > 1, the item is a stack. Crafting decrements the stack. When quantity reaches 0, the row is deleted.

---

## 4. Room Tags

Two approaches:

### Option A: JSON field on Room

Add `tags []string` as a JSON field on the Room schema. Simple, no new schema, but no relational queries (can't do "find all rooms tagged `forge`" easily).

### Option B: Room → Tag M2M edge

Create a M2M edge from Room to the existing Tag entity. Reuses the admin tag management UI. Lets you query "which rooms have forge tags" — useful for map features later.

**Recommendation: Option A for v1 (JSON field), migrate to Option B later.** The `craft` command only needs to check the current room, not search across rooms.

---

## 5. API Endpoints

### 5.1 Recipe CRUD

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/recipes` | List recipes (filterable by world, station_tag) |
| GET | `/api/recipes/:name` | Get recipe detail |
| POST | `/api/recipes` | Create recipe (admin) |
| PUT | `/api/recipes/:name` | Update recipe (admin) |
| DELETE | `/api/recipes/:name` | Delete recipe (admin) |

### 5.2 Craft Endpoint

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/characters/:id/craft` | Execute craft attempt |

Request:
```json
{
    "recipe": "pepperoni_pizza"
}
```

Response (success):
```json
{
    "success": true,
    "outputs": [
        {"template_id": "pepperoni_pizza", "quantity": 1, "instance_id": 1234}
    ]
}
```

Response (failure):
```json
{
    "success": false,
    "error": "missing ingredient: dough",
    "missing_inputs": ["dough"]
}
```

---

## 6. Server-Side Crafting Logic

The `POST /api/characters/:id/craft` handler:

1. **Lookup recipe** — find `CraftingRecipe` by name
2. **Check station** — query room's `tags` field for `required_station_tag` (need room ID from character's currentRoomId)
3. **Check class** — if recipe has `required_class`, verify character has that faction category membership
4. **Check skill level** — if recipe has skill requirement, check character's competency
5. **Check inventory** — for each `input`, verify the character owns `quantity` instances of matching equipment templates
6. **Consume inputs** — decrement quantities, delete empty stacks
7. **Spawn outputs** — create Equipment instances from templates, assign to character's inventory
8. **Return result** — list of new item IDs

Errors are descriptive: "missing ingredient: dough", "need pizza_station in this room".

---

## 7. Client-Side Commands

### `craft <recipe>`

```
> craft pepperoni_pizza
Checking recipe...
Standing at pizza_station ✓
Ingredients: dough x1, sauce x1, cheese x1, pepperoni x1 ✓
Crafting pepperoni pizza... (3s)
Crafted! Pepperoni Pizza added to inventory.
```

### `recipes [search]`

```
> recipes
Available recipes at pizza_station:
  1. pepperoni_pizza — Pepperoni Pizza (requires: dough, sauce, cheese, pepperoni)
  2. cheese_pizza — Cheese Pizza (requires: dough, sauce, cheese)
  3. supreme_pizza — Supreme Pizza (requires: dough, sauce, cheese, pepperoni, mushroom, olive)
```

### `stations`

```
> stations
This room has the following stations: pizza_station
Type 'recipes pizza_station' to see available recipes.
```

---

## 8. The Pizza Chef Class

Using the existing FactionCategory system:

- **FactionCategory:** "class" with `initial_config: true`
- **Faction:** "pizza_chef" with `display_name: "Pizza Chef"`, `member_tags: ["pizza_chef"]`
- The `required_class` field on a recipe matches the faction name ("pizza_chef")

When a player picks the Pizza Chef class during character creation:
1. They get the `pizza_chef` tag
2. Recipes with `required_class: "pizza_chef"` become available
3. Pizza recipes produce weapons that deal damage based on ingredient quality

---

## 9. Pizza as a Weapon

The "Pepperoni Pizza" is an **EquipmentTemplate** with:

```json
{
    "name": "pepperoni_pizza",
    "display_name": "Pepperoni Pizza",
    "item_type": "weapon",
    "slot": "main_hand",
    "damage_dice_count": 1,
    "damage_dice_sides": 4,
    "damage_bonus": 2,
    "damage_type": "bludgeoning",
    "weapon_type": "improvised",
    "is_two_handed": false,
    "stats": {
        "hunger_restored": 15,
        "is_edible": 1
    }
}
```

The template is created through the admin UI, then referenced by the recipe as an output. Nothing special needed — the existing equip/combat system handles it.

---

## 10. Implementation Order

| Phase | What | Depends On |
|-------|------|------------|
| **1** | Room tags (JSON field on Room) + migration | Nothing |
| **2** | `CraftingRecipe` schema + ent codegen | Phase 1 |
| **3** | Recipe CRUD routes + admin UI page | Phase 2 |
| **4** | Server-side craft endpoint (POST /api/characters/:id/craft) | Phase 3 |
| **5** | Client `craft` + `recipes` + `stations` commands | Phase 4 |
| **6** | Item quantity/stacking on Equipment | Phase 1 |
| **7** | Pizza Chef class config + pizza recipes + weapon template | Phase 5 |

---

## 11. Open Questions

1. **Quantity field on Equipment** — breaking change for existing items. All existing rows need `quantity = 1` set. Should this be a migration script or handled in code with a default?

2. **Recipe discovery** — should `recipes` show all recipes the player could craft, or only those they have ingredients for? I lean toward showing all, greyed out if missing ingredients.

3. **Craft time** — should the player sit there for N seconds, or is it instant? Pizza Chef feels instant-fast (3s). A masterwork sword might take 30s. The channel time could be a recipe field or a server-side tick.

4. **Quality system** — should ingredient quality affect output? Dough x1 + Sauce x1 always produces the same pizza. Future: ingredients could have quality/rarity that feeds into output stats.

5. **Multi-output recipes** — e.g., butchering a carcass yields meat + hide + bones. The outputs JSON array handles this. UI needs to support it.

6. **Skill level gating** — should higher skill levels reduce craft time, increase output quality, or unlock better recipes? The `required_skill_level` field gates access; multipliers for time/quality can come later.

---

## 12. Appendix: Existing Patterns Used

| Pattern | Used For |
|---------|----------|
| EquipmentTemplate → Equipment | Recipe outputs |
| FactionCategory (initial_config) | Class selection for Pizza Chef |
| Tag system (CharacterTag) | Class/guild membership gating |
| Room equipment query | Checking room contents for stations |
| JSON fields on ent models | Recipe inputs/outputs, room tags |
| Admin CRUD pages | Recipe management, template creation |
| Effect system | Future: on-craft effects (XP, tags) |
