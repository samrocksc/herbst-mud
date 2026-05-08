-- Migration: Rename skills → abilities (and related tables/columns)
-- Run this BEFORE starting the server with the new schema.
-- The server's ent auto-migration will create the new tables, so this
-- script handles the data migration from old tables to new ones.

BEGIN;

-- 1. Copy data from old skills table to new abilities table
-- (The new abilities table is created by ent auto-migration)
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
-- (The ent auto-migration creates the npc_template_npc_abilities join table)
INSERT INTO npc_template_npc_abilities (npc_template_id, npc_ability_id)
SELECT npc_template_id, npc_skill_id
FROM npc_template_npc_skills
ON CONFLICT DO NOTHING;

-- 5. Update join table references for ability_npc_abilities
INSERT INTO ability_npc_abilities (ability_id, npc_ability_id)
SELECT skill_id, npc_skill_id
FROM skill_npc_skills
ON CONFLICT DO NOTHING;

-- 6. Update faction foreign key references
-- (faction_abilities column replaces faction_skills in abilities table)
-- This is handled by the INSERT in step 1 since faction_skills maps to faction_abilities

-- 7. Drop old tables (after verifying data migration)
-- WARNING: Only run these after verifying all data is migrated successfully!
-- DROP TABLE IF EXISTS skill_npc_skills;
-- DROP TABLE IF EXISTS npc_template_npc_skills;
-- DROP TABLE IF EXISTS npc_skills;
-- DROP TABLE IF EXISTS character_skills;
-- DROP TABLE IF EXISTS skills;

COMMIT;