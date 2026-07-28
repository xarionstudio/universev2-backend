-- Migration Down: Remove seed data from per-category master tables

DELETE FROM master_running_texts;
DELETE FROM master_mess;
DELETE FROM master_locations_ex;
DELETE FROM master_buses;
DELETE FROM master_tempudo;
DELETE FROM master_areas;
DELETE FROM master_eq_classes;
DELETE FROM master_products;
DELETE FROM master_egi_types;