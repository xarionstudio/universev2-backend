-- Migration Up: Add app_settings table and seed all remaining frontend modules

-- 1. APP SETTINGS
CREATE TABLE IF NOT EXISTS app_settings (
    id SERIAL PRIMARY KEY,
    app_name VARCHAR(100) NOT NULL DEFAULT 'UNIVERSE-2',
    app_env VARCHAR(50) NOT NULL DEFAULT 'development',
    company_logo TEXT DEFAULT '',
    theme VARCHAR(20) NOT NULL DEFAULT 'dark',
    lang VARCHAR(10) NOT NULL DEFAULT 'id',
    menu_vis_json TEXT NOT NULL DEFAULT '{"display":true,"roster":true,"employees":true,"ftw":true,"asset":true,"prestasi":true,"master":true,"users":true,"settings":true}',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO app_settings (app_name, app_env, company_logo, theme, lang, menu_vis_json)
VALUES ('UNIVERSE-2', 'development', '', 'dark', 'id', '{"display":true,"roster":true,"employees":true,"ftw":true,"asset":true,"prestasi":true,"master":true,"users":true,"settings":true}');

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
('503264142', CURRENT_DATE, 'N', '',      '',                    '',      '',                    'off');

-- 3. SEED MASTER DATA ENTRIES was in master_entries (dropped by migration 000006).
-- Now they go to per-category tables. Already handled in 000007_seed_master_tables.
-- This section is intentionally left blank since master_entries table no longer exists.

-- 4. SEED FTW LOGS (matches FE initialFtwRows)
INSERT INTO ftw_logs (employee_nik, shift, sleep_minutes, sleep_formatted, status, rest_hours, can_work, send_time, log_date) VALUES
('503264133', 'pagi', 440, '7 j 20 m', 'fit',    0, TRUE,  '05:12', CURRENT_DATE),
('503264134', 'pagi', 405, '6 j 45 m', 'fit',    0, TRUE,  '05:30', CURRENT_DATE),
('503264135', 'pagi', 220, '3 j 40 m', 'pulang', 0, FALSE, '05:01', CURRENT_DATE),
('503264136', 'pagi', 315, '5 j 15 m', 'spare',  1, TRUE,  '05:22', CURRENT_DATE),
('503264137', 'pagi', 390, '6 j 30 m', 'fit',    0, TRUE,  '06:15', CURRENT_DATE),
('503264138', 'pagi', 430, '7 j 10 m', 'fit',    0, TRUE,  '05:40', CURRENT_DATE),
('503264141', 'pagi', 270, '4 j 30 m', 'spare',  2, TRUE,  '05:10', CURRENT_DATE);

-- 5. SEED ROSTER METAS & REVISIONS (matches FE rosterMeta & initialRevisions)
INSERT INTO roster_files (id, label, month_period, dept, filename, total_employees, total_rows, created_by, date_iso, status) VALUES
(1, 'Juli 2026', '2026-07', 'Operation', 'roster_juli_2026_operation.xlsx', 158, '1.364', 'First Angel',    CURRENT_DATE, 'aktif'),
(2, 'Juli 2026', '2026-07', 'Plant',     'roster_juli_2026_plant.xlsx',     46,  '397',   'First Angel',    CURRENT_DATE, 'aktif'),
(3, 'Juli 2026', '2026-07', 'SDI',       'roster_juli_2026_sdi.xlsx',       24,  '207',   'Rahmat Hidayat', CURRENT_DATE, 'aktif'),
(4, 'Juni 2026', '2026-06', 'Operation', 'roster_juni_2026_operation.xlsx', 155, '1.338', 'First Angel',    '2026-06-30', 'arsip');

INSERT INTO roster_revisions (submission_id, employee_nik, what_id, what_en, when_id, when_en, status, by_id, by_en) VALUES
('sub-101', '503264133', 'Ganti shift D ke N', 'Shift change D to N', 'Tgl 15 Jul', '15 Jul', 'pending',  NULL, NULL),
('sub-102', '503264134', 'Cuti tahunan 3 hari', 'Annual leave 3 days', 'Tgl 18-20 Jul', '18-20 Jul', 'approved', 'Disetujui oleh Supervisor', 'Approved by Supervisor');

-- 6. SEED NOTIFICATIONS (matches FE initialNotifs)
-- user_id references users.id (1..5 from migration 000002)
INSERT INTO notifications (id, user_id, tone, text_id, text_en, time_id, time_en, is_read) VALUES
(1, 1, 'warning', '14 revisi absensi menunggu approval', '14 attendance revisions awaiting approval', '5 menit lalu', '5 minutes ago', FALSE),
(2, 1, 'danger',  'Unit DT-114 berstatus Breakdown', 'Unit DT-114 is in Breakdown status', '32 menit lalu', '32 minutes ago', FALSE),
(3, 1, 'success', 'Upload roster Juli berhasil diproses', 'July roster upload processed successfully', '1 jam lalu', '1 hour ago', FALSE),
(4, 1, 'warning', '3 operator Unfit pada shift pagi hari ini', '3 operators unfit on today''s day shift', '2 jam lalu', '2 hours ago', TRUE),
(5, 1, 'danger',  'Display TV Mess 31 terputus dari jaringan', 'TV Mess 31 display disconnected from network', '3 jam lalu', '3 hours ago', TRUE);

-- 7. SEED PRESTASI SCORES & BADGES (matches FE initialPrestasi)
INSERT INTO prestasi_scores (employee_nik, period_days, total_points, rank, streak_days, att_count, sleep_pct, late_count, penalty_count) VALUES
('503264134', 30, 480, 1, 14, 28, 96.5, 0, 0),
('503264133', 30, 452, 2, 10, 27, 94.0, 1, 0),
('503264138', 30, 410, 3,  8, 26, 92.0, 0, 0),
('503264136', 30, 385, 4,  5, 25, 88.5, 1, 0),
('503264137', 30, 360, 5,  4, 24, 85.0, 2, 0);

INSERT INTO prestasi_badges (employee_nik, badge_key) VALUES
('503264134', 'streak14'),
('503264134', 'perfectSleep'),
('503264133', 'neverLate');