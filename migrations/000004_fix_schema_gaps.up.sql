-- Migration Up: Fix schema gaps between DB and frontend requirements
-- Adds missing tables, columns, and fixes prestasi schema

-- 1. Add columns to units_db (matches FE UnitDb type)
ALTER TABLE units_db ADD COLUMN IF NOT EXISTS category VARCHAR(50);
ALTER TABLE units_db ADD COLUMN IF NOT EXISTS upd_date VARCHAR(20) DEFAULT '';
ALTER TABLE units_db ADD COLUMN IF NOT EXISTS upd_by VARCHAR(100) DEFAULT '';

-- 2. Add missing columns to prestasi_scores (matches Go model PrestasiScore)
ALTER TABLE prestasi_scores ADD COLUMN IF NOT EXISTS current_streak INTEGER DEFAULT 0;
ALTER TABLE prestasi_scores ADD COLUMN IF NOT EXISTS sleep_ok_count INTEGER DEFAULT 0;
ALTER TABLE prestasi_scores ADD COLUMN IF NOT EXISTS total_scheduled_days INTEGER DEFAULT 0;
ALTER TABLE prestasi_scores ADD COLUMN IF NOT EXISTS qualified_days INTEGER DEFAULT 0;
ALTER TABLE prestasi_scores ADD COLUMN IF NOT EXISTS cover_days INTEGER DEFAULT 0;
ALTER TABLE prestasi_scores ADD COLUMN IF NOT EXISTS avg_sleep_min INTEGER DEFAULT 0;

-- 3. Create prestasi_history table (matches Go model PrestasiHistoryEntry + FE PrestasiDay)
CREATE TABLE IF NOT EXISTS prestasi_history (
    id SERIAL PRIMARY KEY,
    employee_nik VARCHAR(50) REFERENCES employees(nik) ON DELETE CASCADE,
    record_date DATE NOT NULL,
    period_days INTEGER NOT NULL DEFAULT 30,
    shift_code VARCHAR(20),
    unit_code VARCHAR(50),
    att_status VARCHAR(20),
    clock_in VARCHAR(10),
    is_late BOOLEAN DEFAULT FALSE,
    sleep_min INTEGER DEFAULT 0,
    att_ok BOOLEAN DEFAULT FALSE,
    sleep_ok BOOLEAN DEFAULT FALSE,
    ftw_status VARCHAR(20),
    rest_hours INTEGER DEFAULT 0,
    outcome VARCHAR(30),
    counterpart_nik VARCHAR(50) DEFAULT '',
    counterpart_name VARCHAR(150) DEFAULT '',
    points INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_prestasi_history_nik_date ON prestasi_history(employee_nik, record_date);

-- 4. Update existing prestasi_scores seed data with new columns
UPDATE prestasi_scores SET
    current_streak = streak_days,
    sleep_ok_count = CAST(ROUND(sleep_pct / 100.0 * period_days) AS INTEGER),
    total_scheduled_days = period_days,
    qualified_days = att_count,
    cover_days = 0,
    avg_sleep_min = 420
WHERE current_streak IS NULL;
