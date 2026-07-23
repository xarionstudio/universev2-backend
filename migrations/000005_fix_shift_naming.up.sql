-- Migration Up: Fix shift naming in fleet_allocations to match FE ("pagi" not "siang")

-- 1. Change CHECK constraint on fleet_allocations.shift
ALTER TABLE fleet_allocations DROP CONSTRAINT IF EXISTS fleet_allocations_shift_check;
ALTER TABLE fleet_allocations ADD CONSTRAINT fleet_allocations_shift_check CHECK (shift IN ('pagi', 'malam'));

-- 2. Update existing seed data
UPDATE fleet_allocations SET shift = 'pagi' WHERE shift = 'siang';