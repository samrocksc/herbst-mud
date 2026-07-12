🔵 feat(abilities): class-specific abilities for Ooze Surfers + eligibility fix

## New Class Abilities

Created 4 abilities each for the two Ooze Surfers classes:

### Trash Mage (faction_id=59)
- **Trash Bolt** — INT-scaled damage, 8 mana, 4 cooldown
- **Junk Shield** — damage absorption buff, 3 rounds, 12 mana
- **Putrid Spray** — damage + accuracy debuff, 10 mana
- **Salvage Aura** — mana regen over 3 rounds

### Foot Clank (faction_id=60)
- **Mech Blade Slash** — DEX-scaled slashing damage, 12 stamina
- **Cloak Field** — dodge buff, 2 rounds, 15 stamina
- **Servo Stomp** — damage + stun, 18 stamina
- **System Reboot** — self-heal, 10 stamina

Each ability has effects in ability_effects table with appropriate
scaling stats, durations, and effect types (damage, buff, debuff,
stun, heal).

## Bug Fix: ability eligibility service missing WithFaction()

The characterActiveFactionIDs() query in ability_eligibility.go
did not call .WithFaction() when loading character_factions. This
meant m.Edges.Faction was always nil, so faction IDs were never
added to the active set. All faction-linked abilities showed
eligible=false with reason "not_active_member_of_faction" even
when the character had an active membership.

Fix: Added .WithFaction() to the CharacterFaction query.

## Cleanup

- Deleted duplicated "Classes" category (id=7) and its pizza_chef
  faction (id=39) — redundant with the "class" category (id=10)
- Deleted empty "Professions" category (id=8)
- Fixed smack's equipped ability: was pointing to world 1 Slap
  (id=10), now points to world 2 Slap (id=27)
- Added character_factions membership for smack → trash_mage
  (faction 59, status=active)

## DB State After Changes

Categories:
  - class (id=9, world 1, initial_config=true) — 8 dev classes
  - class (id=10, world 2, initial_config=true) — trash_mage, foot_clank

Class abilities (world 2):
  - ids 40-43: Trash Mage abilities (faction 59)
  - ids 44-47: Foot Clank abilities (faction 60)

Files:
  server/service/ability_eligibility.go  +WithFaction() on CharacterFaction query
  features/player-class-abilities.feature  5 Gherkin scenarios

Verified via browser: Trash Bolt equipped to smack, POST 201 confirmed.
All 4 trash_mage abilities eligible=true, all 4 foot_clank abilities
eligible=false for smack.