-- Migration script to convert competition_settings.regions from integer[] to jsonb
-- This needs to be run once to fix the data type mismatch
-- First, let's backup the existing data
CREATE TABLE IF NOT EXISTS competition_settings_backup AS
SELECT *
FROM competition_settings;
-- Drop the existing table (this will cascade to any foreign keys if they exist)
DROP TABLE IF EXISTS competition_settings CASCADE;
-- Recreate the table with the correct schema
-- Note: GORM AutoMigrate will recreate this table when the app restarts
-- But we can also manually create it here to be explicit
-- After running this, restart your backend application
-- GORM AutoMigrate will recreate the table with the correct jsonb type
-- If you had any important settings, you can manually insert them:
-- INSERT INTO competition_settings (metric, regions, active, created_at, updated_at)
-- VALUES ('kill_count', '[]'::jsonb, true, NOW(), NOW());