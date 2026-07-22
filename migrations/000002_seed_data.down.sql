-- Migration Seed Down: Clean up seed data

DELETE FROM unit_statuses WHERE code IN ('DT5108', 'EX5002', 'DT114');
DELETE FROM units_db WHERE id IN ('udb-1', 'udb-2', 'udb-3');
DELETE FROM employees WHERE nik IN ('503264133', '503264134', '503264138');
DELETE FROM user_roles WHERE user_id IN ('u1', 'u2', 'u3');
DELETE FROM users WHERE id IN ('u1', 'u2', 'u3');
DELETE FROM role_permissions WHERE role_id IN ('r1', 'r2', 'r3');
DELETE FROM roles WHERE id IN ('r1', 'r2', 'r3');
