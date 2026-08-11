-- Migration Down: Revert business rules settings table

-- ============================================================================
-- Business Rules Settings — Remove table
-- ============================================================================

DROP TABLE IF EXISTS business_rules;