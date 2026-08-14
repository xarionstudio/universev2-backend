ALTER TABLE roster_revisions ADD COLUMN IF NOT EXISTS target_date DATE;
CREATE INDEX IF NOT EXISTS idx_roster_revisions_target_date ON roster_revisions(target_date);
