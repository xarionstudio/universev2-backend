-- Migration Up: Rename FTW shift value 'siang' to 'pagi' to match fleet allocation.

-- 1. Fix CHECK constraint on ftw_logs.shift (DROP first so we can update)
ALTER TABLE ftw_logs DROP CONSTRAINT IF EXISTS ftw_logs_shift_check;

-- 2. Update existing FTW logs (day shift value)
UPDATE ftw_logs SET shift = 'pagi' WHERE shift = 'siang';

-- 3. Add new CHECK constraint back with 'pagi' and 'malam'
ALTER TABLE ftw_logs ADD CONSTRAINT ftw_logs_shift_check CHECK (shift IN ('pagi', 'malam'));