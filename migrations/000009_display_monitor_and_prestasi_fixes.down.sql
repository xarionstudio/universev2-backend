-- Migration Down: Revert display monitor multi-fleet + prestasi fixes

-- ============================================================================
-- 1. Prestasi — hapus indeks
-- ============================================================================

DROP INDEX IF EXISTS idx_prestasi_history_unit;

-- ============================================================================
-- 2. Audio Schedule — kembalikan CHECK constraint
-- ============================================================================

ALTER TABLE audio_schedule_displays DROP CONSTRAINT IF EXISTS audio_schedule_displays_display_kind_check;
ALTER TABLE audio_schedule_displays ADD CONSTRAINT audio_schedule_displays_display_kind_check
  CHECK (display_kind IN ('att', 'fleet', 'ftw', 'finger'));

-- ============================================================================
-- 3. Display Monitor — hapus pivot + kolom + kembalikan CHECK
-- ============================================================================

DROP TABLE IF EXISTS display_fleets;

ALTER TABLE display_devices DROP COLUMN IF EXISTS rotate_sec;

ALTER TABLE display_devices DROP CONSTRAINT IF EXISTS display_devices_content_kind_check;
ALTER TABLE display_devices ADD CONSTRAINT display_devices_content_kind_check
  CHECK (content_kind IN ('att', 'fleet', 'ftw', 'finger'));