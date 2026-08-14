-- Migration Down: Revert FTW shift value 'pagi' back to 'siang'.

-- 1. Revert existing data
UPDATE ftw_logs SET shift = 'siang' WHERE shift = 'pagi';

-- 2. Restore old CHECK constraint
ALTER TABLE ftw_logs DROP CONSTRAINT IF EXISTS ftw_logs_shift_check;
ALTER TABLE ftw_logs ADD CONSTRAINT ftw_logs_shift_check CHECK (shift IN ('siang', 'malam'));