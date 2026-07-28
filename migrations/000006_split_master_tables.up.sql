-- Migration Up: Split master_entries into per-category tables with SERIAL PK + code field

-- ============================================================================
-- 1. Create individual master tables with descriptive column names
-- ============================================================================

CREATE TABLE IF NOT EXISTS master_egi_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_products (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_eq_classes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_areas (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_tempudo (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) DEFAULT '',
    pickup_type VARCHAR(100) DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_buses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    egi_type VARCHAR(255) DEFAULT '',
    departure_time VARCHAR(10) DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_locations_ex (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    bus_code VARCHAR(50) DEFAULT '',
    tempudo_code VARCHAR(50) DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_mess (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    block VARCHAR(100) DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS master_running_texts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    target_display VARCHAR(100) DEFAULT '',
    text_color VARCHAR(50) DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 2. Migrate data from master_entries to new tables
-- ============================================================================

INSERT INTO master_egi_types (code, name, is_active, created_at, updated_at)
SELECT id, name, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'egi';

INSERT INTO master_products (code, name, is_active, created_at, updated_at)
SELECT id, name, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'product';

INSERT INTO master_eq_classes (code, name, description, is_active, created_at, updated_at)
SELECT id, name, field_a, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'eqclass';

INSERT INTO master_areas (code, name, category, is_active, created_at, updated_at)
SELECT id, name, field_a, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'area';

INSERT INTO master_tempudo (code, name, location, pickup_type, is_active, created_at, updated_at)
SELECT id, name, field_a, field_b, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'tempudo';

INSERT INTO master_buses (code, name, egi_type, departure_time, is_active, created_at, updated_at)
SELECT id, name, field_a, field_b, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'bus';

INSERT INTO master_locations_ex (code, name, bus_code, tempudo_code, is_active, created_at, updated_at)
SELECT id, name, field_a, field_b, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'lokasiex';

INSERT INTO master_mess (code, name, block, is_active, created_at, updated_at)
SELECT id, name, field_a, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'mess';

INSERT INTO master_running_texts (code, name, target_display, text_color, is_active, created_at, updated_at)
SELECT id, name, field_a, field_b, is_active, created_at, updated_at
FROM master_entries WHERE category_key = 'runtext';

-- ============================================================================
-- 3. Drop old master_entries table (data already migrated)
-- ============================================================================

DROP TABLE IF EXISTS master_entries;

-- ============================================================================
-- 4. Indexes on code for fast lookups
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_master_egi_types_code ON master_egi_types(code);
CREATE INDEX IF NOT EXISTS idx_master_products_code ON master_products(code);
CREATE INDEX IF NOT EXISTS idx_master_eq_classes_code ON master_eq_classes(code);
CREATE INDEX IF NOT EXISTS idx_master_areas_code ON master_areas(code);
CREATE INDEX IF NOT EXISTS idx_master_tempudo_code ON master_tempudo(code);
CREATE INDEX IF NOT EXISTS idx_master_buses_code ON master_buses(code);
CREATE INDEX IF NOT EXISTS idx_master_locations_ex_code ON master_locations_ex(code);
CREATE INDEX IF NOT EXISTS idx_master_mess_code ON master_mess(code);
CREATE INDEX IF NOT EXISTS idx_master_running_texts_code ON master_running_texts(code);