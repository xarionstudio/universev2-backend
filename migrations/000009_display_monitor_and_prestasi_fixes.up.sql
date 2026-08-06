-- Migration Up: Support display monitor multi-fleet + prestasi fixes

-- ============================================================================
-- 1. Display Monitor — multi-fleet + rotateSec
-- ============================================================================

-- Perluas content_kind untuk mendukung "monitor"
ALTER TABLE display_devices DROP CONSTRAINT IF EXISTS display_devices_content_kind_check;
ALTER TABLE display_devices ADD CONSTRAINT display_devices_content_kind_check
  CHECK (content_kind IN ('att', 'fleet', 'ftw', 'finger', 'monitor'));

-- Tambah kolom rotate_sec (durasi per giliran, default 10)
ALTER TABLE display_devices ADD COLUMN IF NOT EXISTS rotate_sec INTEGER DEFAULT 10;

-- Tabel pivot: satu display monitor menampung banyak fleet berurutan
CREATE TABLE IF NOT EXISTS display_fleets (
    display_id INTEGER REFERENCES display_devices(id) ON DELETE CASCADE,
    fleet_id INTEGER REFERENCES fleet_settings(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (display_id, fleet_id)
);

-- ============================================================================
-- 2. Audio Schedule — dukung display monitor
-- ============================================================================

ALTER TABLE audio_schedule_displays DROP CONSTRAINT IF EXISTS audio_schedule_displays_display_kind_check;
ALTER TABLE audio_schedule_displays ADD CONSTRAINT audio_schedule_displays_display_kind_check
  CHECK (display_kind IN ('att', 'fleet', 'ftw', 'finger', 'monitor'));

-- ============================================================================
-- 3. Prestasi — indeks untuk query klasemen per eq class
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_prestasi_history_unit ON prestasi_history(unit_code);

-- ============================================================================
-- 4. Seed display_fleets — hubungkan display monitor ke fleet settings
--    (hanya jika display monitor & fleet settings sudah ada)
-- ============================================================================

-- Seed DSP-M01: 12 fleet (3 halaman x 4 fleet)
INSERT INTO display_fleets (display_id, fleet_id, sort_order)
SELECT d.id, f.id, row_number() OVER (ORDER BY f.id) - 1
FROM display_devices d
CROSS JOIN fleet_settings f
WHERE d.code = 'DSP-M01'
  AND f.is_active = TRUE
  AND f.id IN (
    SELECT id FROM fleet_settings ORDER BY id LIMIT 12
  )
ON CONFLICT (display_id, fleet_id) DO NOTHING;

-- Seed DSP-M02: 6 fleet (2 halaman)
INSERT INTO display_fleets (display_id, fleet_id, sort_order)
SELECT d.id, f.id, row_number() OVER (ORDER BY f.id) - 1
FROM display_devices d
CROSS JOIN fleet_settings f
WHERE d.code = 'DSP-M02'
  AND f.is_active = TRUE
  AND f.id IN (
    SELECT id FROM fleet_settings ORDER BY id OFFSET 12 LIMIT 6
  )
ON CONFLICT (display_id, fleet_id) DO NOTHING;
