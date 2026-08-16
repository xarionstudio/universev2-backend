-- Migration Down: hapus unit hasil seed 000016 — hanya yang persis
-- snapshot dummy yang dimasukkan migrasi ini (lokasi kosong, tidak standby,
-- tidak breakdown, kosong upd). Unit hasil import manual / seed lama TIDAK
-- boleh terhapus oleh down migration (kondisinya berbeda dari snapshot ini).
DELETE FROM units_db
WHERE is_active = TRUE
  AND is_standby = FALSE
  AND is_breakdown = FALSE
  AND (location IS NULL OR location = '')
  AND (upd_date IS NULL OR upd_date = '')
  AND (upd_by IS NULL OR upd_by = '')
  AND category IS NOT NULL AND category <> '';