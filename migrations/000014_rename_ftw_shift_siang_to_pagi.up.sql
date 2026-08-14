-- Migration Up: Rename FTW shift value 'siang' to 'pagi' to match fleet allocation.

-- 1. Update existing FTW logs (day shift value)
UPDATE ftw_logs SET shift = 'pagi' WHERE shift = 'siang';

-- 2. Fix CHECK constraint on ftw_logs.shift
ALTER TABLE ftw_logs DROP CONSTRAINT IF EXISTS ftw_logs_shift_check;
ALTER TABLE ftw_logs ADD CONSTRAINT ftw_logs_shift_check CHECK (shift IN ('pagi', 'malam'));