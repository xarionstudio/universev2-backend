-- Migration Down: Revert shift naming change

-- 1. Revert existing data
UPDATE fleet_allocations SET shift = 'siang' WHERE shift = 'pagi';

-- 2. Restore old CHECK constraint
ALTER TABLE fleet_allocations DROP CONSTRAINT IF EXISTS fleet_allocations_shift_check;
ALTER TABLE fleet_allocations ADD CONSTRAINT fleet_allocations_shift_check CHECK (shift IN ('siang', 'malam'));