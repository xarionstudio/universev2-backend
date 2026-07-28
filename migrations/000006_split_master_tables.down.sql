-- Migration Down: Drop per-category master tables (data preserved in master_entries)

DROP TABLE IF EXISTS master_running_texts;
DROP TABLE IF EXISTS master_mess;
DROP TABLE IF EXISTS master_locations_ex;
DROP TABLE IF EXISTS master_buses;
DROP TABLE IF EXISTS master_tempudo;
DROP TABLE IF EXISTS master_areas;
DROP TABLE IF EXISTS master_eq_classes;
DROP TABLE IF EXISTS master_products;
DROP TABLE IF EXISTS master_egi_types;