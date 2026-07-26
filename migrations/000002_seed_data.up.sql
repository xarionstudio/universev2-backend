-- Migration Seed: Initial Master & Auth Data

-- ============================================================================
-- 1. ROLES & PERMISSIONS
-- ============================================================================

INSERT INTO roles (id, name, description, is_locked) VALUES
('r1', 'Superadmin', 'Semua modul + user & role', TRUE),
('r2', 'Admin', 'Operasional harian — roster, fit to work, fleet, master', FALSE),
('r3', 'Viewer', 'Hanya lihat — tanpa aksi ubah', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Superadmin permissions (full manage on all 10 modules)
INSERT INTO role_permissions (role_id, module_name, permission_level) VALUES
('r1', 'dashboard', 'manage'),
('r1', 'display',   'manage'),
('r1', 'employees', 'manage'),
('r1', 'roster',    'manage'),
('r1', 'ftw',       'manage'),
('r1', 'asset',     'manage'),
('r1', 'prestasi',  'manage'),
('r1', 'master',    'manage'),
('r1', 'users',     'manage'),
('r1', 'settings',  'manage'),

-- Admin permissions
('r2', 'dashboard', 'view'),
('r2', 'display',   'manage'),
('r2', 'employees', 'manage'),
('r2', 'roster',    'manage'),
('r2', 'ftw',       'manage'),
('r2', 'asset',     'manage'),
('r2', 'prestasi',  'view'),
('r2', 'master',    'manage'),
('r2', 'users',     'none'),
('r2', 'settings',  'none'),

-- Viewer permissions
('r3', 'dashboard', 'view'),
('r3', 'display',   'none'),
('r3', 'employees', 'view'),
('r3', 'roster',    'view'),
('r3', 'ftw',       'view'),
('r3', 'asset',     'view'),
('r3', 'prestasi',  'view'),
('r3', 'master',    'view'),
('r3', 'users',     'none'),
('r3', 'settings',  'none')
ON CONFLICT (role_id, module_name) DO NOTHING;

-- ============================================================================
-- 2. SEED INITIAL USERS (matches FE initialUmUsers)
-- Password default: "admin123" (sha256 representation)
-- ============================================================================

INSERT INTO users (id, email, name, nik, password_hash, password_salt, is_active) VALUES
('u1', 'angel@unggul.co.id',     'First Angel Paustine', '503264133', 'fd7555bb94a4534a601551a2635ed8f4b631d92ca1ff16f73a9f28ab73fd4cf9', 'saltsuperadmin', TRUE),
('u2', 'rahmat.h@unggul.co.id',  'Rahmat Hidayat',       '503264134', '67d094a72b5ec32000a0867a0bd0313feada1a4133e85460f6073d4e4919ebf2', 'saltadmin',      TRUE),
('u3', 'dewi.l@unggul.co.id',    'Dewi Lestari',         '503264138', '2559dd25634a0d370ded4af91110691c31ab00ae51616060e3d6b2ef1f150778', 'saltadmin2',     TRUE),
('u4', 'clinic@unggul.co.id',    'Klinik Viewer',        NULL,        '0da3712e37a2b0ea15cca65e93cf91bbc93efac82d3414cbf749c99dbbb0b5b9', 'saltviewer1',    TRUE),
('u5', 'budi.plant@unggul.co.id','Hendra Gunawan',       '503264143', 'd4c150faba808bc315df9a36d02cf4d8f03df4780ac469bbb2b51af1f45f487b', 'saltviewer2',    FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id) VALUES
('u1', 'r1'),
('u2', 'r2'),
('u3', 'r2'),
('u4', 'r3'),
('u5', 'r3')
ON CONFLICT (user_id, role_id) DO NOTHING;

-- ============================================================================
-- 3. SEED EMPLOYEES (matches FE employees.ts & users)
-- ============================================================================

INSERT INTO employees (nik, name, dept, pos, simper, simper_exp, status, company, equip_type, join_date, exp_date, license_type, mcu_status, blood_type, bpjs_no, mess_name, room_no, phone, emergency_contact) VALUES
('503264133', 'First Angel Paustine', 'Operation', 'Operator Dump Truck',   'DT Kelas A',   '2027-03-14', 'aktif', 'PT Unggul Dinamika Utama', 'Dump Truck 777D',     '2022-03-14', '2027-03-14', 'SIM B2 Umum', 'Fit — 12 Feb 2026', 'O+', '0001234567890', 'Mess 31 — Blok C', 'C-214', '+62 812-5501-2233', 'Maria Paustine — +62 811-4400-9876'),
('503264134', 'Rahmat Hidayat',       'SDI',       'Foreman Operation',      'Exca Kelas B', '2026-11-20', 'aktif', 'PT Unggul Dinamika Utama', 'Excavator CAT 390F', '2021-05-10', '2026-11-20', 'SIM A',       'Fit — 20 Jan 2026', 'A+', NULL,           'Mess 31 — Blok A', 'A-102', '+62 813-4422-1100', 'Siti Hidayat — +62 812-3311-0099'),
('503264135', 'Budi Santoso',         'HRGA',      'Staf HR',               NULL,           NULL,         'aktif', 'PT Unggul Dinamika Utama', NULL,                 '2023-02-01', NULL,         'SIM C',       'Fit — 10 Jan 2026', 'B-', NULL,           'Mess 31 — Blok D', 'D-101', '+62 814-1234-5678', 'Rina Santoso — +62 812-9988-0011'),
('503264136', 'Siti Nurhaliza',       'Operation', 'Operator Dump Truck',   'DT Kelas A',   '2027-08-15', 'aktif', 'PT Unggul Dinamika Utama', 'Dump Truck 777D',     '2022-08-15', '2027-08-15', 'SIM B2 Umum', 'Fit — 01 Mar 2026', 'AB+',NULL,           'Mess 31 — Blok C', 'C-210', '+62 819-8877-6655', 'Hendra N — +62 812-3344-5566'),
('503264137', 'Andi Prasetyo',        'Plant',     'Mekanik Senior',        NULL,           NULL,         'aktif', 'PT Unggul Dinamika Utama', NULL,                 '2020-06-01', NULL,         'SIM C',       'Fit — 15 Feb 2026', 'O-', NULL,           'Mess 31 — Blok B', 'B-305', '+62 811-5544-3322', 'Wati Prasetyo — +62 817-2211-0099'),
('503264138', 'Dewi Lestari',         'HRGA',      'Staff Roster & HR',     NULL,           '2028-01-01', 'aktif', 'PT Unggul Dinamika Utama', NULL,                 '2023-01-15', '2028-01-01', 'SIM C',       'Fit — 05 Jan 2026', 'B+', NULL,           'Mess 31 — Blok D', 'D-305', '+62 815-9988-7766', 'Budi Lestari — +62 815-1122-3344'),
('503264139', 'Joko Widodo S.',       'Operation', 'Operator Excavator',    'Exca Kelas A', '2027-05-10', 'aktif', 'PT Unggul Dinamika Utama', 'Excavator CAT 390F', '2021-05-10', '2027-05-10', 'SIM B2 Umum', 'Fit — 22 Jan 2026', 'A-', NULL,           'Mess 31 — Blok A', 'A-215', '+62 813-5566-7788', 'Sari Widodo — +62 818-1122-3344'),
('503264140', 'Rina Marlina',         'HRGA',      'Staf Administrasi',     NULL,           NULL,         'aktif', 'PT Unggul Dinamika Utama', NULL,                 '2022-09-01', NULL,         'SIM C',       'Fit — 30 Jan 2026', 'O+', NULL,           NULL,               NULL,    '+62 822-3344-5566', 'Budi Marlina — +62 822-1122-0099'),
('503264141', 'Agus Salim',           'Plant',     'Operator Grader',       'GR Kelas B',   '2026-12-01', 'aktif', 'PT Unggul Dinamika Utama', 'Grader 140M',        '2021-12-01', '2026-12-01', 'SIM B2 Umum', 'Fit — 08 Feb 2026', 'B+', NULL,           'Mess 31 — Blok B', 'B-201', '+62 817-7788-9900', 'Nani Salim — +62 819-3344-2211'),
('503264142', 'Maya Sari',            'Operation', 'Operator Dump Truck',   'DT Kelas A',   '2027-10-20', 'aktif', 'PT Unggul Dinamika Utama', 'Dump Truck 777D',     '2022-10-20', '2027-10-20', 'SIM B2 Umum', 'Fit — 19 Feb 2026', 'A+', NULL,           'Mess 31 — Blok C', 'C-312', '+62 821-9900-1122', 'Rudi Sari — +62 821-8877-6655'),
('503264143', 'Hendra Gunawan',       'Plant',     'Teknisi Alat Berat',    NULL,           NULL,         'aktif', 'PT Unggul Dinamika Utama', NULL,                 '2022-03-10', NULL,         'SIM C',       'Fit — 25 Jan 2026', 'AB-',NULL,           'Mess 31 — Blok B', 'B-412', '+62 812-0011-2233', 'Sri Gunawan — +62 812-5544-3322')
ON CONFLICT (nik) DO NOTHING;

-- Employee competencies
INSERT INTO employee_competencies (employee_nik, class_name, simper_no, expiry_date) VALUES
('503264133', 'DT', 'DT Kelas A', '2027-03-14'),
('503264134', 'EX', 'Exca Kelas B', '2026-11-20'),
('503264139', 'EX', 'Exca Kelas A', '2027-05-10'),
('503264141', 'GR', 'GR Kelas B', '2026-12-01'),
('503264142', 'DT', 'DT Kelas A', '2027-10-20')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 4. SEED UNITS DB & STATUSES (matches FE units-db.ts / unit-status.ts)
-- ============================================================================

INSERT INTO units_db (id, code, egi, product, class_name, work_area, is_active, is_standby, is_breakdown, location) VALUES
('udb-1',  'DT5108', 'DT777', 'Caterpillar', 'DT', 'Panel East Tengah', TRUE,  FALSE, FALSE, 'Panel East Tengah'),
('udb-2',  'EX5002', 'EX700', 'Hitachi',     'EX', 'Panel East Tengah', TRUE,  FALSE, FALSE, 'Panel East Tengah'),
('udb-3',  'DT114',  'DT777', 'Caterpillar', 'DT', 'Pit Service',       TRUE,  FALSE, TRUE,  'Workshop Plant'),
('udb-4',  'DT5112', 'DT777', 'Caterpillar', 'DT', 'Panel East Tengah', TRUE,  FALSE, FALSE, 'Panel East Tengah'),
('udb-5',  'DT5111', 'DT777', 'Caterpillar', 'DT', 'Panel East Tengah', TRUE,  FALSE, FALSE, 'Panel East Tengah'),
('udb-6',  'EX7001', 'EX700', 'Caterpillar', 'EX', 'Workshop Plant',    TRUE,  TRUE,  FALSE, 'Workshop Plant'),
('udb-7',  'EX7007', 'EX700', 'Caterpillar', 'EX', 'Kantor Operation',  TRUE,  FALSE, FALSE, 'Kantor Operation'),
('udb-8',  'EX8001', 'EX800', 'Hitachi',     'EX', 'KM 31',             TRUE,  FALSE, FALSE, 'KM 31')
ON CONFLICT (code) DO NOTHING;

INSERT INTO unit_statuses (code, type_info, status, location, updated_note) VALUES
('DT5108', 'DT777 · Caterpillar', 'ready',     'Panel East Tengah', 'Checklist harian ok'),
('EX5002', 'EX700 · Hitachi',     'ready',     'Panel East Tengah', 'Checklist harian ok'),
('DT114',  'DT777 · Caterpillar', 'breakdown', 'Workshop Plant',    '04:15 — dilaporkan rusak, menunggu perbaikan'),
('DT5112', 'DT777 · Caterpillar', 'ready',     'Panel East Tengah', 'Checklist harian ok'),
('DT5111', 'DT777 · Caterpillar', 'ready',     'Panel East Tengah', 'Checklist harian ok'),
('EX7001', 'EX700 · Caterpillar', 'standby',   'Workshop Plant',    'Cadangan shift pagi'),
('EX7007', 'EX700 · Caterpillar', 'ready',     'Kantor Operation',  'Checklist harian ok'),
('EX8001', 'EX800 · Hitachi',     'ready',     'KM 31',             'Checklist harian ok')
ON CONFLICT (code) DO NOTHING;

-- Unit status histories (hist field — 4-tuple: when, what, why, status)
INSERT INTO unit_status_histories (unit_code, hist_when, hist_what, hist_why, hist_status) VALUES
('DT5108', '13 Jul 05:15', 'Ready',     'Checklist harian ok',                    'ready'),
('EX5002', '13 Jul 05:10', 'Ready',     'Checklist harian ok',                    'ready'),
('DT114',  '12 Jul 04:15', 'Breakdown', 'Dilaporkan rusak — masuk antrean workshop','breakdown'),
('DT114',  '05 Jul 06:00', 'Ready',     'Checklist harian ok',                    'ready'),
('EX7001', '12 Jul 06:00', 'Standby',   'Cadangan shift pagi',                    'standby'),
('EX7001', '10 Jul 05:00', 'Ready',     'Checklist harian ok',                    'ready');

-- ============================================================================
-- 5. SEED FLEET SETTINGS (matches FE initialFleets)
-- ============================================================================

INSERT INTO fleet_settings (id, digger_code, location, bus_code, is_active) VALUES
('fl-EX5002', 'EX5002', 'Panel East Tengah', 'UDBU001', TRUE),
('fl-EX7001', 'EX7001', 'Workshop Plant',    'UDBU002', TRUE),
('fl-EX7007', 'EX7007', 'Kantor Operation',  'UDBU003', TRUE),
('fl-EX8001', 'EX8001', 'KM 31',            'UDBU004', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO fleet_setting_units (fleet_setting_id, unit_code) VALUES
('fl-EX5002', 'DT5108'),
('fl-EX5002', 'DT5112'),
('fl-EX5002', 'DT5111'),
('fl-EX7001', 'DT114')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 6. SEED AUDIO SCHEDULES (matches FE initialAudios)
-- ============================================================================

INSERT INTO audio_schedules (id, title, trigger_time, frequency, file_name, is_active) VALUES
('au1', 'Pengumuman P5M',    '05:45', 'harian', 'p5m_reminder.mp3',  TRUE),
('au2', 'Alarm fatigue check','13:00', 'perjam', 'fatigue_alert.mp3', TRUE),
('au3', 'Pergantian shift',  '17:45', 'harian', 'shift_change.mp3',  FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO audio_schedule_displays (audio_id, display_kind) VALUES
('au1', 'att'),
('au1', 'fleet'),
('au2', 'ftw'),
('au3', 'att'),
('au3', 'fleet'),
('au3', 'ftw'),
('au3', 'finger')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 7. SEED DISPLAY DEVICES (matches FE initialDspAtt + initialDspFleet)
-- ============================================================================

INSERT INTO display_devices (id, name, location, content_kind, fleet_id, running_text, is_online, last_heartbeat, is_active) VALUES
('DSP-A01', 'TV Gate Utara',    'Gate utara',   'att',   NULL,         'Utamakan keselamatan — patuhi batas kecepatan 40 km/jam di jalan hauling.', TRUE,  '3 dtk lalu',  TRUE),
('DSP-A02', 'TV Mess 31',       'Lobi mess',    'att',   NULL,         'Wajib P2H sebelum mengoperasikan unit.',                                    TRUE,  '12 dtk lalu', TRUE),
('DSP-A03', 'TV Ruang P5M',     'Kantor SDI',   'ftw',   NULL,         'Rapat P5M setiap pergantian shift di front masing-masing.',                 FALSE, '26 mnt lalu', TRUE),
('DSP-A04', 'TV Pos Fingerprint','Gate selatan', 'finger',NULL,         'Wajib P2H sebelum mengoperasikan unit.',                                    TRUE,  '5 dtk lalu',  FALSE),
('DSP-F01', 'Fleet EX7001',     'Workshop Plant','fleet', 'fl-EX7001',  'Wajib P2H sebelum mengoperasikan unit.',                                    TRUE,  '2 dtk lalu',  TRUE),
('DSP-F02', 'Fleet EX7007',     'Kantor Operation','fleet','fl-EX7007', 'Utamakan keselamatan — patuhi batas kecepatan 40 km/jam di jalan hauling.', TRUE,  '8 dtk lalu',  TRUE),
('DSP-F03', 'Fleet EX8001',     'KM 31',        'fleet', 'fl-EX8001',  'Musim hujan: waspadai jalan licin di ramp Pit Tempudo.',                    FALSE, '1 j lalu',    TRUE)
ON CONFLICT (id) DO NOTHING;
