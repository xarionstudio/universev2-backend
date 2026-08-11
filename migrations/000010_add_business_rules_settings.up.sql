-- Migration Up: Add business rules settings table

-- ============================================================================
-- Business Rules Settings — JSON-based configuration for frontend constants
-- ============================================================================

CREATE TABLE IF NOT EXISTS business_rules (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50) NOT NULL UNIQUE,
    rules JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(100)
);

-- Seed initial values from current hardcoded constants
INSERT INTO business_rules (category, rules, updated_by) VALUES
(
    'prestasi',
    '{
        "pts_base": 10,
        "pts_ontime": 2,
        "pts_sleep": 3,
        "pts_streak_step": 2,
        "pts_streak_cap": 10,
        "pts_cover": 5,
        "sleep_min_great": 420,
        "period_days": {
            "week": 7,
            "month": 30,
            "quarter": 90
        }
    }',
    'system'
),
(
    'ftw',
    '{
        "sleep_fit_min": 330,
        "sleep_spare_1h_min": 300,
        "sleep_spare_2h_min": 240
    }',
    'system'
),
(
    'fleet',
    '{
        "max_units": 13
    }',
    'system'
),
(
    'auth',
    '{
        "password_min_length": 8
    }',
    'system'
),
(
    'weather',
    '{
        "refresh_interval_ms": 900000
    }',
    'system'
)
ON CONFLICT (category) DO NOTHING;

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_business_rules_category ON business_rules(category);