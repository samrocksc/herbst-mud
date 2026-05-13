-- Migration: Add main_type column to quests table
-- This adds the quest type categorization field

ALTER TABLE quests
ADD COLUMN IF NOT EXISTS main_type VARCHAR(32) DEFAULT 'general';
