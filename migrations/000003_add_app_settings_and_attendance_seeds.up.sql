-- Migration Up: Add app_settings table and seed all remaining frontend modules

-- 1. APP SETTINGS
CREATE TABLE IF NOT EXISTS app_settings (
    id VARCHAR(50) PRIMARY KEY DEFAULT 'default',
    app_name VARCHAR(100) NOT NULL DEFAULT 'UNIVERSE-2',
    app_env VARCHAR(50) NOT NULL DEFAULT 'development',
    company_logo TEXT DEFAULT '',
    theme VARCHAR(20) NOT NULL DEFAULT 'dark',
    lang VARCHAR(10) NOT NULL DEFAULT 'id',
    menu_vis_json TEXT NOT NULL DEFAULT '{"display":true,"roster":true,"employees":true,"ftw":true,"asset":true,"prestasi":true,"master":true,"users":true,"settings":true}',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO app_settings (id, app_name, app_env, company_logo, theme, lang, menu_vis_json)
VALUES ('default', 'UNIVERSE-2', 'development', '', 'dark', 'id', '{"display":true,"roster":true,"employees":true,"ftw":true,"asset":true,"prestasi":true,"master":true,"users":true,"settings":true}')
ON CONFLICT (id) DO NOTHING;

-- 2. SEED ATTENDANCE LOGS
INSERT INTO attendance_logs (employee_nik, attendance_date, shift_code, check_in, check_in_machine, check_out, check_out_machine, status) VALUES
('503264133', CURRENT_DATE, 'D', '05:45', 'FP-02 · Gate utara',   '17:32', 'FP-02 · Gate utara',   'hadir'),
('503264134', CURRENT_DATE, 'D', '07:02', 'FP-01 · Office',      '',      '',                    'hadir'),
('503264135', CURRENT_DATE, 'D', '',      '',                    '',      '',                    'unfit'),
('503264136', CURRENT_DATE, 'D', '05:51', 'FP-03 · Gate selatan', '',      '',                    'hadir'),
('503264137', CURRENT_DATE, 'D', '06:31', 'FP-04 · Workshop',    '',      '',                    'terlambat'),
('503264138', CURRENT_DATE, 'D', '06:58', 'FP-01 · Office',      '',      '',                    'hadir'),
('503264139', CURRENT_DATE, 'D', '',      '',                    '',      '',                    'belum'),
('503264140', CURRENT_DATE, 'CR', '',     '',                    '',      '',                    'off'),
('503264141', CURRENT_DATE, 'D', '',      '',                    '',      '',                    'unfit'),
('503264142', CURRENT_DATE, 'N', '',      '',                    '',      '',                    'off')
ON CONFLICT DO NOTHING;

-- 3. SEED MASTER DATA ENTRIES (matches FE mdInit())
INSERT INTO master_entries (id, category_key, name, field_a, field_b, is_active) VALUES
('egi-0',     'egi',      'DT777',                '',                        '',         TRUE),
('egi-1',     'egi',      'EX700',                '',                        '',         TRUE),
('egi-2',     'egi',      'EX800',                '',                        '',         TRUE),
('product-0', 'product',  'Caterpillar',          '',                        '',         TRUE),
('product-1', 'product',  'Hitachi',              '',                        '',         TRUE),
('product-2', 'product',  'Komatsu',              '',                        '',         TRUE),
('area-0',    'area',     'Panel East Tengah',   'Mining',                  '',         TRUE),
('area-1',    'area',     'Workshop Plant',       'Non-Mining',              '',         TRUE),
('area-2',    'area',     'Kantor Operation',     'Non-Mining',              '',         TRUE),
('tempudo-0', 'tempudo',  'TP-01',                'Workshop',                'Pickup',   TRUE),
('tempudo-1', 'tempudo',  'TP-02',                'Panel East Tengah',       'Pickup & Drop', TRUE),
('mess-0',    'mess',     'Mess 31',              'Blok A',                  '',         TRUE),
('mess-1',    'mess',     'Mess 31',              'Blok C',                  '',         TRUE),
('mess-2',    'mess',     'Mess 31',              'Blok D',                  '',         TRUE),
('mess-3',    'mess',     'Mess KM 12',           'Blok B',                  '',         TRUE),
('runtext-0', 'runtext',  'Utamakan keselamatan — patuhi batas kecepatan 40 km/jam di jalan hauling.', 'Semua kiosk', 'Cyan', TRUE),
('runtext-1', 'runtext',  'Wajib P2H sebelum mengoperasikan unit.', 'Display Fleet', 'Oranye', TRUE),
('runtext-2', 'runtext',  'Rapat P5M setiap pergantian shift di front masing-masing.', 'Display Attendance', 'Putih', TRUE),
('runtext-3', 'runtext',  'Musim hujan: waspadai jalan licin di ramp Pit Tempudo.', 'Semua kiosk', 'Merah', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 4. SEED FTW LOGS (matches FE initialFtwRows)
INSERT INTO ftw_logs (employee_nik, shift, sleep_minutes, sleep_formatted, status, rest_hours, can_work, send_time, log_date) VALUES
('503264133', 'siang', 440, '7 j 20 m', 'fit',    0, TRUE,  '05:12', CURRENT_DATE),
('503264134', 'siang', 405, '6 j 45 m', 'fit',    0, TRUE,  '05:30', CURRENT_DATE),
('503264135', 'siang', 220, '3 j 40 m', 'pulang', 0, FALSE, '05:01', CURRENT_DATE),
('503264136', 'siang', 315, '5 j 15 m', 'spare',  1, TRUE,  '05:22', CURRENT_DATE),
('503264137', 'siang', 390, '6 j 30 m', 'fit',    0, TRUE,  '06:15', CURRENT_DATE),
('503264138', 'siang', 430, '7 j 10 m', 'fit',    0, TRUE,  '05:40', CURRENT_DATE),
('503264141', 'siang', 270, '4 j 30 m', 'spare',  2, TRUE,  '05:10', CURRENT_DATE)
ON CONFLICT DO NOTHING;

-- 5. SEED ROSTER METAS & REVISIONS (matches FE rosterMeta & initialRevisions)
INSERT INTO roster_files (id, label, month_period, dept, filename, total_employees, total_rows, created_by, date_iso, status) VALUES
('jul',       'Juli 2026', '2026-07', 'Operation', 'roster_juli_2026_operation.xlsx', 158, '1.364', 'First Angel',    CURRENT_DATE, 'aktif'),
('jul-plant', 'Juli 2026', '2026-07', 'Plant',     'roster_juli_2026_plant.xlsx',     46,  '397',   'First Angel',    CURRENT_DATE, 'aktif'),
('jul-sdi',   'Juli 2026', '2026-07', 'SDI',       'roster_juli_2026_sdi.xlsx',       24,  '207',   'Rahmat Hidayat', CURRENT_DATE, 'aktif'),
('jun',       'Juni 2026', '2026-06', 'Operation', 'roster_juni_2026_operation.xlsx', 155, '1.338', 'First Angel',    '2026-06-30', 'arsip')
ON CONFLICT (id) DO NOTHING;

INSERT INTO roster_revisions (submission_id, employee_nik, what_id, what_en, when_id, when_en, status, by_id, by_en) VALUES
('sub-101', '503264133', 'Ganti shift D ke N', 'Shift change D to N', 'Tgl 15 Jul', '15 Jul', 'pending',  NULL, NULL),
('sub-102', '503264134', 'Cuti tahunan 3 hari', 'Annual leave 3 days', 'Tgl 18-20 Jul', '18-20 Jul', 'approved', 'Disetujui oleh Supervisor', 'Approved by Supervisor')
ON CONFLICT DO NOTHING;

-- 6. SEED NOTIFICATIONS (matches FE initialNotifs)
INSERT INTO notifications (id, user_id, tone, text_id, text_en, time_id, time_en, is_read) VALUES
('n1', 'u1', 'warning', '14 revisi absensi menunggu approval', '14 attendance revisions awaiting approval', '5 menit lalu', '5 minutes ago', FALSE),
('n2', 'u1', 'danger',  'Unit DT-114 berstatus Breakdown', 'Unit DT-114 is in Breakdown status', '32 menit lalu', '32 minutes ago', FALSE),
('n3', 'u1', 'success', 'Upload roster Juli berhasil diproses', 'July roster upload processed successfully', '1 jam lalu', '1 hour ago', FALSE),
('n4', 'u1', 'warning', '3 operator Unfit pada shift pagi hari ini', '3 operators unfit on today''s day shift', '2 jam lalu', '2 hours ago', TRUE),
('n5', 'u1', 'danger',  'Display TV Mess 31 terputus dari jaringan', 'TV Mess 31 display disconnected from network', '3 jam lalu', '3 hours ago', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 7. SEED PRESTASI SCORES & BADGES (matches FE initialPrestasi)
INSERT INTO prestasi_scores (employee_nik, period_days, total_points, rank, streak_days, att_count, sleep_pct, late_count, penalty_count) VALUES
('503264134', 30, 480, 1, 14, 28, 96.5, 0, 0),
('503264133', 30, 452, 2, 10, 27, 94.0, 1, 0),
('503264138', 30, 410, 3,  8, 26, 92.0, 0, 0),
('503264136', 30, 385, 4,  5, 25, 88.5, 1, 0),
('503264137', 30, 360, 5,  4, 24, 85.0, 2, 0)
ON CONFLICT DO NOTHING;

INSERT INTO prestasi_badges (employee_nik, badge_key) VALUES
('503264134', 'streak14'),
('503264134', 'perfectSleep'),
('503264133', 'neverLate')
ON CONFLICT DO NOTHING;
