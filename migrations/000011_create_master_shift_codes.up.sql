-- Create master_shift_codes table
CREATE TABLE IF NOT EXISTS master_shift_codes (
    id SERIAL PRIMARY KEY,
    group_id VARCHAR(50),
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for category lookups
CREATE INDEX IF NOT EXISTS idx_master_shift_codes_code ON master_shift_codes(code);
CREATE INDEX IF NOT EXISTS idx_master_shift_codes_active ON master_shift_codes(is_active);

-- Seed default shift codes
INSERT INTO master_shift_codes (group_id, code, name, name_en, sort_order, is_active) VALUES
-- Shift & kehadiran
('Shift & kehadiran', 'D', 'Day shift', 'Day shift', 1, TRUE),
('Shift & kehadiran', 'N', 'Night shift', 'Night shift', 2, TRUE),
('Shift & kehadiran', 'R', 'Reguler', 'Regular', 3, TRUE),
('Shift & kehadiran', 'STB', 'Standby', 'Standby', 4, TRUE),
('Shift & kehadiran', 'OFF', 'OFF', 'OFF', 5, TRUE),

-- Cuti & izin
('Cuti & izin', 'CR', 'Cuti roster', 'Roster leave', 10, TRUE),
('Cuti & izin', 'AL', 'Annual leave', 'Annual leave', 11, TRUE),
('Cuti & izin', 'LWP', 'Izin dengan upah', 'Paid leave', 12, TRUE),
('Cuti & izin', 'LWOP', 'Izin tanpa upah', 'Unpaid leave', 13, TRUE),
('Cuti & izin', 'PH', 'Public holiday', 'Public holiday', 14, TRUE),
('Cuti & izin', 'PHD', 'Public holiday siang', 'Public holiday (day)', 15, TRUE),

-- Sakit & ketidakhadiran
('Sakit & ketidakhadiran', 'S', 'Sakit', 'Sick', 20, TRUE),
('Sakit & ketidakhadiran', 'A', 'Alpha', 'Alpha / no notice', 21, TRUE),

-- Medis & karantina
('Medis & karantina', 'MCU', 'Medical check up', 'Medical check up', 30, TRUE),
('Medis & karantina', 'MCR', 'Reguler MCU', 'Regular MCU', 31, TRUE),
('Medis & karantina', 'MCUF', 'Follow up MCU', 'MCU follow-up', 32, TRUE),
('Medis & karantina', 'ISM', 'Isolasi mandiri', 'Self-isolation', 33, TRUE),
('Medis & karantina', 'OBC', 'Observasi COVID', 'COVID observation', 34, TRUE),
('Medis & karantina', 'KRT', 'Karantina', 'Quarantine', 35, TRUE),

-- Tugas & training
('Tugas & training', 'TGS', 'Tugas', 'Assignment', 40, TRUE),
('Tugas & training', 'DNS', 'Dinas', 'Official duty', 41, TRUE),
('Tugas & training', 'TRV', 'Travel', 'Travel', 42, TRUE),
('Tugas & training', 'TR', 'Training di luar site', 'Off-site training', 43, TRUE),
('Tugas & training', 'TRS', 'Training onsite', 'On-site training', 44, TRUE),
('Tugas & training', 'IN', 'Induksi', 'Induction', 45, TRUE),

-- Status kepegawaian
('Status kepegawaian', 'TERM', 'Termination', 'Termination', 50, TRUE),
('Status kepegawaian', 'EOC', 'Kontrak berakhir', 'Contract ended', 51, TRUE),
('Status kepegawaian', 'RSG', 'Resign', 'Resign', 52, TRUE)
ON CONFLICT (code) DO NOTHING;