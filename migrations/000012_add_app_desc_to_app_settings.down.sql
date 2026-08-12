-- Migration Down: Remove app_desc column from app_settings
ALTER TABLE app_settings DROP COLUMN IF EXISTS app_desc;