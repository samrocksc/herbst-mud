-- Migration: Add world_id validation constraints
-- Fixes invalid world_ids and adds DB-level referential integrity
-- Run this manually via psql after deploying the schema changes.

-- Step 1: Fix any remaining invalid world_ids across all content tables
UPDATE abilities SET world_id = '1' WHERE world_id IS NULL OR world_id = '' OR world_id = 'default';
UPDATE crafting_recipes SET world_id = '1' WHERE world_id IS NULL OR world_id = '' OR world_id = 'default';
UPDATE dialog_nodes SET world_id = '1' WHERE world_id IS NULL OR world_id = '' OR world_id = 'default';
UPDATE equipment_templates SET world_id = '1' WHERE world_id IS NULL OR world_id = '' OR world_id = 'default';
UPDATE factions SET world_id = '1' WHERE world_id IS NULL OR world_id = '' OR world_id = 'default';
UPDATE npc_templates SET world_id = '1' WHERE world_id IS NULL OR world_id = '' OR world_id = 'default';
UPDATE quests SET world_id = '1' WHERE world_id IS NULL OR world_id = '' OR world_id = 'default';
UPDATE rooms SET world_id = '1' WHERE world_id IS NULL OR world_id = '' OR world_id = 'default';

-- Step 2: Create validation function
-- Ensures world_id matches an existing worlds.id (cast bigint to text for comparison)
CREATE OR REPLACE FUNCTION validate_world_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.world_id IS NULL OR NEW.world_id = '' THEN
        NEW.world_id := '1';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM worlds WHERE id::text = NEW.world_id) THEN
        RAISE EXCEPTION 'Invalid world_id: %. Must match a valid worlds.id', NEW.world_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Step 3: Apply trigger to all content tables with world_id
-- Skips applogs (optional world context) and worlds itself
DO $$
DECLARE
    tname TEXT;
BEGIN
    FOR tname IN
        SELECT table_name
        FROM information_schema.columns
        WHERE column_name = 'world_id'
          AND table_schema = 'public'
          AND table_name != 'worlds'
          AND table_name != 'applogs'
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS validate_world_id_trigger ON %I', tname);
        EXECUTE format('CREATE TRIGGER validate_world_id_trigger BEFORE INSERT OR UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION validate_world_id()', tname);
    END LOOP;
END $$;
