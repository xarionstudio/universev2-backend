DROP INDEX IF EXISTS idx_roster_revisions_target_date;
ALTER TABLE roster_revisions DROP COLUMN IF EXISTS target_date;
