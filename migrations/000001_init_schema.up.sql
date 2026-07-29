-- Migration Up: Initialize UniverseV2 PostgreSQL Schema
-- All primary keys use SERIAL (auto-increment numeric)

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- 1. AUTH & RBAC (Users, Roles, Permissions)
-- ============================================================================

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    nik VARCHAR(50),
    password_hash VARCHAR(255) NOT NULL,
    password_salt VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    password_changed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    -- module_name matches FE UmModule: dashboard|display|employees|roster|ftw|asset|prestasi|master|users|settings
    module_name VARCHAR(50) NOT NULL,
    -- permission_level matches FE UmPerm: none|view|manage
    permission_level VARCHAR(20) NOT NULL CHECK (permission_level IN ('none', 'view', 'manage')),
    PRIMARY KEY (role_id, module_name)
);

-- ============================================================================
-- 2. KARYAWAN & MASTER DATA
-- ============================================================================

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    nik VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    dept VARCHAR(100) NOT NULL,
    pos VARCHAR(100) NOT NULL,
    simper VARCHAR(100),
    simper_exp DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'aktif' CHECK (status IN ('aktif', 'cuti', 'nonaktif')),
    company VARCHAR(150) NOT NULL,
    equip_type VARCHAR(100),
    join_date DATE,
    exp_date DATE,
    license_type VARCHAR(100),
    mcu_status VARCHAR(150),
    med_history TEXT,
    blood_type VARCHAR(10),
    bpjs_no VARCHAR(50),
    mess_name VARCHAR(100),
    room_no VARCHAR(50),
    phone VARCHAR(50),
    emergency_contact VARCHAR(255),
    photo_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employee_competencies (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees(id) ON DELETE CASCADE,
    class_name VARCHAR(100) NOT NULL,
    simper_no VARCHAR(100) NOT NULL,
    expiry_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS units_db (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    egi VARCHAR(100),
    product VARCHAR(100),
    class_name VARCHAR(100),
    work_area VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    is_standby BOOLEAN DEFAULT FALSE,
    is_breakdown BOOLEAN DEFAULT FALSE,
    location VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 3. ROSTER & ABSENSI
-- ============================================================================

CREATE TABLE IF NOT EXISTS roster_files (
    id SERIAL PRIMARY KEY,
    label VARCHAR(100) NOT NULL,
    month_period VARCHAR(7) NOT NULL, -- Format YYYY-MM
    dept VARCHAR(100) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    total_employees INTEGER DEFAULT 0,
    total_rows VARCHAR(50),
    created_by VARCHAR(150),
    date_iso DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'aktif' CHECK (status IN ('aktif', 'arsip')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roster_schedules (
    id SERIAL PRIMARY KEY,
    roster_file_id INTEGER REFERENCES roster_files(id) ON DELETE CASCADE,
    employee_nik VARCHAR(50) REFERENCES employees(nik) ON DELETE CASCADE,
    schedule_date DATE NOT NULL,
    -- shift_code matches FE legend codes: D, N, R, STB, OFF, CR, AL, LWP, LWOP, S, A, MCU, etc.
    shift_code VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS roster_revisions (
    id SERIAL PRIMARY KEY,
    -- submission_id = sid in FE ApRow
    submission_id VARCHAR(50) NOT NULL,
    employee_nik VARCHAR(50) REFERENCES employees(nik) ON DELETE CASCADE,
    -- what_id/what_en = whatId/whatEn in FE ApRow
    what_id TEXT NOT NULL,
    what_en TEXT NOT NULL,
    -- when_id/when_en = whenId/whenEn in FE ApRow
    when_id TEXT NOT NULL,
    when_en TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    -- by_id/by_en = byId/byEn in FE ApRow
    by_id TEXT,
    by_en TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS attendance_logs (
    id SERIAL PRIMARY KEY,
    employee_nik VARCHAR(50) REFERENCES employees(nik) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    shift_code VARCHAR(20),
    check_in VARCHAR(10),         -- FE: in
    check_in_machine VARCHAR(100), -- FE: inM
    check_out VARCHAR(10),        -- FE: out
    check_out_machine VARCHAR(100),-- FE: outM
    -- status matches FE AttStatus: hadir|terlambat|belum|unfit|off
    status VARCHAR(20) CHECK (status IN ('hadir', 'terlambat', 'belum', 'unfit', 'off')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 4. FIT TO WORK (FTW)
-- ============================================================================

CREATE TABLE IF NOT EXISTS ftw_logs (
    id SERIAL PRIMARY KEY,
    employee_nik VARCHAR(50) REFERENCES employees(nik) ON DELETE CASCADE,
    -- shift matches FE FtwRecord.shift: siang|malam
    shift VARCHAR(20) CHECK (shift IN ('siang', 'malam')),
    -- sleep_minutes = sleepMin in FE
    sleep_minutes INTEGER,
    -- sleep_formatted = sleep in FE (display string "7 j 20 m")
    sleep_formatted VARCHAR(20),
    -- status matches FE FtwStatus: fit|spare|pulang|belum
    status VARCHAR(20) CHECK (status IN ('fit', 'spare', 'pulang', 'belum')),
    -- rest_hours = restHours in FE
    rest_hours INTEGER DEFAULT 0,
    can_work BOOLEAN DEFAULT TRUE,
    send_time VARCHAR(20),
    log_date DATE NOT NULL,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 5. ASSETS & FLEET MANAGEMENT
-- ============================================================================

CREATE TABLE IF NOT EXISTS unit_statuses (
    id SERIAL PRIMARY KEY,
    unit_code VARCHAR(50) UNIQUE NOT NULL REFERENCES units_db(code) ON DELETE CASCADE,
    type_info VARCHAR(150),
    -- status matches FE UnitStatus: ready|breakdown|standby
    status VARCHAR(20) NOT NULL CHECK (status IN ('ready', 'breakdown', 'standby')),
    location VARCHAR(150),
    updated_note TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Stores history as 4-field rows matching FE UnitHist = [when, what, why, status]
CREATE TABLE IF NOT EXISTS unit_status_histories (
    id SERIAL PRIMARY KEY,
    unit_code VARCHAR(50) REFERENCES unit_statuses(unit_code) ON DELETE CASCADE,
    -- "when" column: timestamp label (e.g. "12 Jul 04:15")
    hist_when VARCHAR(50) NOT NULL,
    -- "what" column: status label (e.g. "Breakdown", "Ready")
    hist_what VARCHAR(100) NOT NULL,
    -- "why" column: reason text
    hist_why TEXT,
    -- "status" column: UnitStatus value for dot color
    hist_status VARCHAR(20) NOT NULL CHECK (hist_status IN ('ready', 'breakdown', 'standby')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fleet_settings (
    id SERIAL PRIMARY KEY,
    digger_code VARCHAR(50) NOT NULL,
    location VARCHAR(150) NOT NULL,
    bus_code VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fleet_setting_units (
    fleet_setting_id INTEGER REFERENCES fleet_settings(id) ON DELETE CASCADE,
    unit_code VARCHAR(50) NOT NULL,
    PRIMARY KEY (fleet_setting_id, unit_code)
);

CREATE TABLE IF NOT EXISTS fleet_allocations (
    id SERIAL PRIMARY KEY,
    alloc_date DATE NOT NULL,
    -- shift matches FE: pagi|malam
    shift VARCHAR(20) NOT NULL CHECK (shift IN ('pagi', 'malam')),
    fleet_id INTEGER REFERENCES fleet_settings(id) ON DELETE CASCADE,
    digger_code VARCHAR(50) NOT NULL,
    location VARCHAR(150),
    bus_code VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fleet_allocation_operators (
    allocation_id INTEGER REFERENCES fleet_allocations(id) ON DELETE CASCADE,
    unit_code VARCHAR(50) NOT NULL,
    operator_nik VARCHAR(50) REFERENCES employees(nik) ON DELETE SET NULL,
    PRIMARY KEY (allocation_id, unit_code)
);

-- ============================================================================
-- 6. PRESTASI & GAMIFIKASI
-- ============================================================================

CREATE TABLE IF NOT EXISTS prestasi_scores (
    id SERIAL PRIMARY KEY,
    employee_nik VARCHAR(50) REFERENCES employees(nik) ON DELETE CASCADE,
    period_days INTEGER NOT NULL DEFAULT 30,
    total_points INTEGER NOT NULL DEFAULT 0,
    rank INTEGER,
    streak_days INTEGER DEFAULT 0,   -- FE: streak
    att_count INTEGER DEFAULT 0,     -- FE: attCount
    sleep_pct NUMERIC(5,2) DEFAULT 0.0, -- FE: sleepPct
    late_count INTEGER DEFAULT 0,
    penalty_count INTEGER DEFAULT 0,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS prestasi_badges (
    id SERIAL PRIMARY KEY,
    employee_nik VARCHAR(50) REFERENCES employees(nik) ON DELETE CASCADE,
    -- badge_key matches FE badge strings: streak14, perfectSleep, neverLate, etc.
    badge_key VARCHAR(50) NOT NULL,
    awarded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 7. NOTIFIKASI, SETTINGS & DISPLAY TV
-- ============================================================================

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    -- tone matches FE Notif.tone: info|success|warning|danger
    tone VARCHAR(20) CHECK (tone IN ('info', 'success', 'warning', 'danger')),
    text_id TEXT NOT NULL,   -- FE: textId
    text_en TEXT NOT NULL,   -- FE: textEn
    time_id VARCHAR(50),     -- FE: timeId (relative string "5 menit lalu")
    time_en VARCHAR(50),     -- FE: timeEn
    is_read BOOLEAN DEFAULT FALSE, -- FE: read
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS audio_schedules (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    -- trigger_time → FE: when (e.g. "05:45")
    trigger_time VARCHAR(10) NOT NULL,
    -- frequency matches FE Audio.freq: sekali|harian|perjam|per30
    frequency VARCHAR(20) NOT NULL CHECK (frequency IN ('sekali', 'harian', 'perjam', 'per30')),
    -- file_name → FE: file
    file_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Junction table: which display types each audio plays on
-- FE: Audio.displays = ("att" | "fleet" | "ftw" | "finger")[]
CREATE TABLE IF NOT EXISTS audio_schedule_displays (
    audio_id INTEGER REFERENCES audio_schedules(id) ON DELETE CASCADE,
    display_kind VARCHAR(20) NOT NULL CHECK (display_kind IN ('att', 'fleet', 'ftw', 'finger')),
    PRIMARY KEY (audio_id, display_kind)
);

CREATE TABLE IF NOT EXISTS display_devices (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    location VARCHAR(150) NOT NULL,
    -- content_kind matches FE DisplayKind: att|fleet|ftw|finger
    content_kind VARCHAR(20) NOT NULL CHECK (content_kind IN ('att', 'fleet', 'ftw', 'finger')),
    -- fleet_id for fleet-type displays (FE: Display.fleetId)
    fleet_id INTEGER REFERENCES fleet_settings(id) ON DELETE SET NULL,
    -- running_text → FE: runtext
    running_text TEXT,
    -- is_online → FE: online
    is_online BOOLEAN DEFAULT TRUE,
    -- last_heartbeat → FE: hb (relative string "3 dtk lalu")
    last_heartbeat VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- INDEXES FOR OPTIMAL QUERY PERFORMANCE
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_employees_nik ON employees(nik);
CREATE INDEX IF NOT EXISTS idx_employees_dept ON employees(dept);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees(status);
CREATE INDEX IF NOT EXISTS idx_roster_schedules_date ON roster_schedules(schedule_date);
CREATE INDEX IF NOT EXISTS idx_attendance_logs_date ON attendance_logs(attendance_date);
CREATE INDEX IF NOT EXISTS idx_ftw_logs_nik_date ON ftw_logs(employee_nik, log_date);
CREATE INDEX IF NOT EXISTS idx_unit_statuses_status ON unit_statuses(status);
CREATE INDEX IF NOT EXISTS idx_unit_status_histories_code ON unit_status_histories(unit_code);
CREATE INDEX IF NOT EXISTS idx_fleet_allocations_date_shift ON fleet_allocations(alloc_date, shift);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, is_read);