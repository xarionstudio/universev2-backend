-- Migration Up: Add app_desc column to app_settings
ALTER TABLE app_settings ADD COLUMN IF NOT EXISTS app_desc TEXT DEFAULT '';