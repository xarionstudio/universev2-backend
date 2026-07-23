-- Migration Down: Revert schema gaps fix

-- 1. Remove columns from units_db
ALTER TABLE units_db DROP COLUMN IF EXISTS category;
ALTER TABLE units_db DROP COLUMN IF EXISTS upd_date;
ALTER TABLE units_db DROP COLUMN IF EXISTS upd_by;

-- 2. Drop prestasi_history table
DROP TABLE IF EXISTS prestasi_history;

-- 3. Remove added columns from prestasi_scores
ALTER TABLE prestasi_scores DROP COLUMN IF EXISTS current_streak;
ALTER TABLE prestasi_scores DROP COLUMN IF EXISTS sleep_ok_count;
ALTER TABLE prestasi_scores DROP COLUMN IF EXISTS total_scheduled_days;
ALTER TABLE prestasi_scores DROP COLUMN IF EXISTS qualified_days;
ALTER TABLE prestasi_scores DROP COLUMN IF EXISTS cover_days;
ALTER TABLE prestasi_scores DROP COLUMN IF EXISTS avg_sleep_min;
