-- Migration: Rename skills → abilities, merge talents → abilities
-- Run this BEFORE starting the server with the new schema.
-- The server's ent auto-migration will create the new tables, so this
-- script handles the data migration from old tables to new ones.

BEGIN;

-- ============================================================
-- PART 1: Copy skills → abilities (rename migration)
-- ============================================================

-- 1. Copy data from old skills table to new abilities table
INSERT INTO abilities (
    id, name, description, ability_type, cost, cooldown, requirements,
    effect_type, effect_value, effect_duration,
    scaling_stat, scaling_percent_per_point,
    mana_cost, stamina_cost, hp_cost,
    slug, required_tag, ability_class, proc_chance, proc_event, cooldown_seconds,
    faction_abilities
)
SELECT
    id, name, description, skill_type, cost, cooldown, requirements,
    effect_type, effect_value, effect_duration,
    scaling_stat, scaling_percent_per_point,
    mana_cost, stamina_cost, hp_cost,
    slug, required_tag, skill_class, proc_chance, proc_event, cooldown_seconds,
    faction_skills
FROM skills
ON CONFLICT (id) DO NOTHING;

-- 2. Copy data from old character_skills to new character_abilities
INSERT INTO character_abilities (
    id, slot, ability_characters, character_abilities
)
SELECT
    id, slot, skill_characters, character_skills
FROM character_skills
ON CONFLICT (id) DO NOTHING;

-- 3. Copy data from old npc_skills to new npc_abilities
INSERT INTO npc_abilities (id, slot)
SELECT id, slot
FROM npc_skills
ON CONFLICT (id) DO NOTHING;

-- 4. Update join table references for npc_abilities
INSERT INTO npc_template_npc_abilities (npc_template_id, npc_ability_id)
SELECT npc_template_id, npc_skill_id
FROM npc_template_npc_skills
ON CONFLICT DO NOTHING;

-- 5. Update join table references for ability_npc_abilities
INSERT INTO ability_npc_abilities (ability_id, npc_ability_id)
SELECT skill_id, npc_skill_id
FROM skill_npc_skills
ON CONFLICT DO NOTHING;

-- ============================================================
-- PART 2: Merge talents → abilities (talent absorption)
-- ============================================================

-- 6. Copy talent rows into abilities as passive abilities
-- Use IDs starting at 1000 to avoid collisions with existing abilities
INSERT INTO abilities (
    id, name, description, ability_type, cost, cooldown, requirements,
    effect_type, effect_value, effect_duration,
    scaling_stat, scaling_percent_per_point,
    mana_cost, stamina_cost, hp_cost,
    slug, required_tag, ability_class, proc_chance, proc_event, cooldown_seconds
)
SELECT
    id + 1000,  -- offset to avoid ID collisions
    name, description, 'passive',  -- talents become passive abilities
    0, cooldown, requirements,
    effect_type, effect_value, effect_duration,
    '', 0,  -- no scaling_stat, no scaling_percent
    mana_cost, stamina_cost, 0,  -- no hp_cost in talents
    '', '',  -- no slug, no required_tag
    'passive',  -- ability_class = passive for former talents
    0, '',  -- no proc_chance, no proc_event
    0  -- no cooldown_seconds
FROM talents
ON CONFLICT (id) DO NOTHING;

-- 7. Copy character_talents into character_abilities
-- Map old talent IDs to new ability IDs (offset by 1000)
INSERT INTO character_abilities (
    slot, ability_characters, character_abilities
)
SELECT
    ct.slot,
    ct.talent_characters + 1000,  -- remap talent ID to ability ID
    ct.character_talents
FROM character_talents ct
WHERE ct.talent_characters IS NOT NULL
ON CONFLICT DO NOTHING;

-- 8. Drop old talent tables (after verifying data migration)
-- WARNING: Only run these after verifying all data is migrated successfully!
-- DROP TABLE IF EXISTS available_talents;
-- DROP TABLE IF EXISTS character_talents;
-- DROP TABLE IF EXISTS talents;

-- ============================================================
-- PART 3: Cleanup (run after verification)
-- ============================================================

-- DROP TABLE IF EXISTS skill_npc_skills;
-- DROP TABLE IF EXISTS npc_template_npc_skills;
-- DROP TABLE IF EXISTS npc_skills;
-- DROP TABLE IF EXISTS character_skills;
-- DROP TABLE IF EXISTS skills;

COMMIT;