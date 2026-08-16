-- Migration Down: kosongkan kembali kolom category (kondisi sebelum migrasi ini).
-- Data kategori unit bukanlah data transaksional, aman untuk di-reset.
ALTER TABLE units_db ADD COLUMN IF NOT EXISTS category VARCHAR(50);
UPDATE units_db SET category = NULL;
