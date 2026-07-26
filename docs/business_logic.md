# UniverseV2 - Business Logic & System Documentation

Dokumentasi ini menjelaskan logika bisnis, fitur, dan alur sistem UniverseV2 yang terdiri dari backend (Go/Fiber) dan frontend (Next.js).

---

## 1. System Overview

### Arsitektur
- **Backend**: Go + Fiber v3 + GORM + PostgreSQL
- **Frontend**: Next.js + React + TypeScript
- **Authentication**: JWT dengan format hash FE (`HashPasswordFE`) dan legacy hash
- **Authorization**: RBAC (Role-Based Access Control) per module dengan permission level `none | view | manage`

### Konsep Dasar
- **Module (UmModule)**: dashboard | display | employees | roster | ftw | asset | prestasi | master | users | settings
- **Permission (UmPerm)**: none | view | manage
- **Role**: Default role `r3` untuk user baru, role `r1` untuk superadmin

---

## 2. Authentication & Authorization

### Fitur
- Login/Logout dengan email dan password
- Registrasi user baru (default role: `r3`)
- Refresh token untuk maintain session
- RBAC middleware proteksi setiap module

### Logika Bisnis
1. **Login**:
   - Validasi email format dan password tidak kosong
   - Cek user aktif (`is_active = true`)
   - Verifikasi password dengan dua format hash:
     - FE format: `HashPasswordFE(password, salt)`
     - Legacy format: `HashPasswordLegacy(password, salt)`
   - Generate JWT token dengan claims: `user_id`, `email`, `roles`
   - Load permissions untuk roles user dari `role_permissions`

2. **RBAC**:
   - Setiap module memiliki permission `view` dan `manage`
   - Handler dilindungi oleh middleware `RequirePermission(module, level)`
   - Contoh: `employees` module butuh permission `view` untuk list, `manage` untuk create/update/delete

---

## 3. Employees Module

### Fitur
- Manajemen data karyawan (NIK, nama, departemen, posisi, status)
- Upload foto karyawan
- Management kompetensi (class_name, simper_no, expiry_date)

### Logika Bisnis
1. **Validasi NIK**: Harus tepat 9 digit
2. **Unik NIK**: Tidak boleh ada 2 karyawan dengan NIK sama
3. **Status Karyawan**: `aktif | cuti | nonaktif`
4. **Kompetensi**: Setiap karyawan bisa punya multiple competencies dengan masa berlaku
5. **Foto**: Upload file photo, disimpan sebagai URL

---

## 4. Fit To Work (FTW)

### Fitur
- Cek kondisi karyawan sebelum kerja (sleep hours, rest hours)
- Log FTW harian per shift (siang/malam)
- History FTW untuk analisis

### Logika Bisnis
1. **Submit FTW**:
   - Input: NIK, shift, sleep_minutes, rest_hours, can_work
   - Sleep format: `7 j 20 m` (jam dan menit)
   - Validasi: NIK harus terdaftar sebagai karyawan aktif
   - `can_work` dikirim dari frontend berdasarkan perhitungan business rule

2. **Status FTW**:
   - `fit`: Bisa bekerja
   - `spare`: Cadangan
   - `pulang`: Dikirim pulang
   - `belum`: Belum dicek

3. **Business Rule** (`FTW Eligibility Logic`):
   - **Operator boleh bekerja (`fit`) jika**:  
     - Memiliki data hadir (`attendance_logs`) ** DAN **  
     - Sleep duration **minimal 5 jam 30 menit (330 menit)**
   - **Kategori `spare` (tidak boleh langsung kerja)**:
     - Sleep 4–4 jam 59 menit: wajib istirahat **2 jam** sebelum boleh bekerja
     - Sleep 5–5 jam 29 menit: wajib istirahat **1 jam** sebelum boleh bekerja
   - **Kategori `pulang`**: Dikirim pulang karena tidak memenuhi syarat
   - **Kategori `belum`**: Belum dicek FTW

4. **Backend Implementation Note**:
   - Saat ini backend hanya menyimpan log FTW (`SubmitLog`)
   - Perhitungan `can_work` dan status `fit/spare/pulang` dilakukan di frontend.
   - Backend hanya menerima nilai `can_work` dari request tanpa verifikasi ulang.

---

## 5. Roster Management

### Fitur
- Upload file roster Excel (per departemen, per bulan)
- Generate jadwal kerja karyawan
- Management revisi jadwal
- Approval workflow revisi
- View absensi berdasarkan roster

### Logika Bisnis
1. **Upload Roster**:
   - File Excel di-upload dan diparse
   - Generate `roster_files` record dengan metadata
   - Setiap baris Excel menjadi `roster_schedules` untuk setiap karyawan
   
2. **Shift Codes**:
   - `D`: Day shift (pagi)
   - `N`: Night shift (malam)
   - `R`: Rest day
   - `STB`: Standby
   - `OFF`: Off
   - `CR`: Consecutive Rest
   - `AL`: Annual Leave
   - `LWP`: Leave Without Pay
   - `LWOP`: Leave Without Pay
   - `S`: Sakit
   - `A`: Absent
   - `MCU`: Medical Check Up

3. **Revisions**:
   - Karyawan/supervisor bisa submit revisi jadwal
   - Approval workflow: `pending → approved | rejected`
   - Support approval dengan catatan

4. **Roster Attendance**:
   - Cross-reference antara `roster_schedules` dan `attendance_logs`
   - Hitung kehadiran vs jadwal

---

## 6. Attendance

### Fitur
- Record check-in/check-out dengan mesin fingerprint
- View attendance hari ini, by date, by range
- Status hadir: `hadir | terlambat | belum | unfit | off`

### Logika Bisnis
1. **Check-In/Out**:
   - Input: NIK, waktu, mesin (FP-01, FP-02, dst)
   - Status **selalu diset `hadir`** saat check-in/check-out.
   - Tidak ada auto-deteksi `terlambat`, `unfit`, `off`, atau `belum` di backend.
   - Tidak ada integrasi dengan FTW untuk validasi attendance.

2. **Backend Implementation Note**:
   - `RecordCheckIn` / `RecordCheckOut` hanya mencatat waktu dan mesin.
   - Semua logika status (terlambat, unfit, dll) **dilakukan di frontend**.

---

## 7. Fleet / Assets Management

### Fitur
- Manajemen fleet setting (digger code → unit assignment)
- Auto allocation unit ke shift
- Status unit: `ready | breakdown | standby`
- History status unit
- Unit DB (master data unit)

### Logika Bisnis
1. **Fleet Setting**:
   - Hubungan: 1 digger code → multiple units
   - Contoh: `EX7001` → `DT114`, `DT115`
   - Setiap fleet setting punya location dan bus code

2. **Auto Allocation**:
   - Algoritma assign unit ke shift berdasarkan:
     - Unit status harus `ready`
     - Tidak ada conflict jadwal
     - Prioritaskan unit yang sudah standby lama

3. **Unit Status**:
   - `ready`: Siap pakai
   - `breakdown`: Rusak, perlu perbaikan
   - `standby`: Cadangan
   - Auto-update status saat breakdown report

4. **Unit History**:
   - Track semua perubahan status
   - Format: `[when, what, why, status]`
   - `when`: timestamp label (e.g., "12 Jul 04:15")
   - `what`: status label (e.g., "Breakdown")
   - `why`: alasan perubahan
   - `status`: untuk warna dot di UI

---

## 8. Prestasi (Performance)

### Fitur
- Leaderboard karyawan berdasarkan poin
- History performa harian
- Badge dan streak tracking

### Logika Bisnis
1. **Score Calculation**:
   - `total_points`: Akumulasi poin 30 hari
   - Poin来源:
     - Attendance: +10 poin hadir, -5 poin terlambat
     - Sleep quality: +5 poin jika sleep > 6 jam
     - Streak: bonus poin untuk konsistensi
   
2. **Streak**:
   - `current_streak`: Hari berturut-turut hadir tanpa terlambat
   - Reset jika ada keterlambatan atau absent

3. **History**:
   - `prestasi_history`: Record harian untuk setiap karyawan
   - Fields: attendance_status, clock_in, sleep_min, ftw_status, points
   - Digenerate nightly batch job

---

## 9. Master Data

### Fitur
- Master data untuk dropdown di frontend
- Categories: dept, pos, simper, mess, egi, product, area, tempudo, runtext, etc.

### Logika Bisnis
1. **Category Management**:
   - Setiap kategori adalah collection of entries
   - Entry aktif/non-aktif (`is_active`)
   - Support field_a dan field_b untuk metadata tambahan

2. **Initialization**:
   - Seeded di migration `000003_add_app_settings_and_attendance_seeds.up.sql`
   - Mirror dengan FE `mdInit()` function

---

## 10. Users & Roles Management

### Fitur
- CRUD users dengan role assignment
- CRUD roles dengan permissions
- Password management (change password)

### Logika Bisnis
1. **User Creation**:
   - Default password strength: min 8 char, letters + numbers, max 72 char
   - Auto-generate salt untuk password hashing
   - Default role: `r3` (operator)

2. **Role Permissions**:
   - Many-to-many: `role_permissions` table
   - Format: `(role_id, module_name, permission_level)`
   - `permission_level`: `none | view | manage`

3. **Password Update**:
   - Verify old password sebelum update
   - Generate new salt dan hash
   - Support fallback dari legacy hash ke FE hash

---

## 11. Notifications

### Fitur
- List notifications untuk user
- Mark read / mark all read
- Notification types: warning, danger, success, info

### Logika Bisnis
1. **Trigger Notifications**:
   - Unit breakdown → danger notification
   - Attendance revision pending → warning notification
   - Successful upload → success notification

2. **Read Status**:
   - `is_read`: boolean flag per notification
   - `MarkAllRead`: update semua notification untuk user_id
   - Frontend real-time fetch notifications

---

## 12. Settings & Displays

### Fitur
- App settings (app_name, theme, lang, menu visibility)
- Audio schedules (pengumuman otomatis)
- Display devices management
- Display heartbeat/monitoring

### Logika Bisnis
1. **App Settings**:
   - Single record dengan id = `default`
   - `menu_vis_json`: JSON object visibility setiap module
   - Theme: `dark | light`
   - Lang: `id | en`

2. **Audio Schedules**:
   - Trigger berbasis waktu (`trigger_time`)
   - Frequency: `harian | perjam`
   - Mapping ke display devices via `audio_schedule_displays`
   - Frontend playlist mode untuk.display devices

3. **Display Devices**:
   - Types: `att | fleet | ftw | finger`
   - Content kind sesuai module FE
   - Heartbeat: frontend kirim ping periodically
   - `is_online`: sync dengan last heartbeat

---

## 13. Fingerprint Devices

### Fitur
- Monitoring status fingerprint sensor devices

### Logika Bisnis
1. **Device Status**:
   - Config: `FINGERPRINT_PORT` (e.g., `/dev/ttyUSB0`), `FINGERPRINT_BAUD` (115200)
   - Frontend display devices dengan `content_kind = 'finger'`
   - Status: `connected | disconnected`
   - Heartbeat tracking

2. **Integration**:
   - Attendance check-in/out menggunakan fingerprint
   - Data dari device: NIK + timestamp
   - Auto-sync ke `attendance_logs`

---

## 14. Dashboard

### Fitur
- Ringkasan data untuk operator
- Metrics: attendance today, unit status, notifications, upcoming shifts

### Logika Bisnis
1. **Summary Metrics**:
   - Attendance rate hari ini
   - Jumlah unit breakdown
   - Notification unread count
   - Karyawan unfit untuk shift berikutnya

2. **Data Aggregation**:
   - Real-time dari `attendance_logs`, `unit_statuses`, `notifications`
   - Pre-calculated untuk performa

---

## 15. Weather Integration

### Fitur
- Cuaca current untuk lokasi operasional

### Logika Bisnis
1. **Proxy API**:
   - Backend proxy ke weather API eksternal
   - Config: koordinat lokasi (lat, lon)
   - Default: `-0.5021, 117.1536` (Samarinda, Kalimantan Timur)

2. **Frontend Usage**:
   - Display di dashboard dan fleet display
   - Impact untuk decision: Operasi atau tidak

---

## 16. Frontend Modules Mapping

| Frontend Module | Backend Routes | Key Entities |
|-----------------|----------------|--------------|
| Dashboard | `/api/weather/current`, aggregated data | attendance_logs, unit_statuses, notifications |
| Displays | `/api/settings/displays`, `/api/displays/:id/heartbeat` | display_devices |
| Employees | `/api/employees/*` | employees, employee_competencies |
| Fit To Work | `/api/ftw/*` | ftw_logs |
| Roster | `/api/rosters/*` | roster_files, roster_schedules, roster_revisions |
| Attendance | `/api/attendance/*` | attendance_logs |
| Asset/Fleet | `/api/fleet/*`, `/api/units/*` | fleet_settings, unit_statuses, unit_status_histories, units_db |
| Prestasi | `/api/prestasi/*` | prestasi_scores, prestasi_badges, prestasi_history |
| Master | `/api/master/*` | master_entries |
| Users | `/api/users/*` | users |
| Roles | `/api/roles/*` | roles, role_permissions |
| Notifications | `/api/notifications/*` | notifications |
| Settings | `/api/settings/*`, `/api/settings/audio`, `/api/settings/displays` | app_settings, audio_schedules, display_devices |
| Fingerprint | `/api/fingerprint/devices` | display_devices (kind=finger) |

---

## 17. Key Business Rules Summary

1. **Password Policy**: Min 8 chars, letters + numbers, max 72 chars
2. **NIK Validation**: Exactly 9 digits, unique per employee
3. **Attendance Treshold**: Terlambat > 15 menit dari shift start
4. **FTW Eligibility**: 
   - **fit**: Attendance hadir ** dan ** sleep >= 5 jam 30 menit
   - **spare**: Sleep 4–4j59m → istirahat 2 jam; Sleep 5–5j29m → istirahat 1 jam
   - **pulang**: Tidak memenuhi syarat
   - Catatan: Saat ini perhitungan dilakukan di frontend, backend hanya menyimpan log
5. **Unit Status**: Hanya unit `ready` yang bisa di-allocation
6. **Roster Upload**: Harus ada kolom NIK, nama, dan shift codes
7. **Display Heartbeat**: Frontend kirim heartbeat setiap 30 detik
8. **JWT Expiration**: Default 24 jam, refresh token 7 hari
9. **RBAC**: Public hanya auth endpoints, protected butuh Bearer token + permission
10. **Role Locking**: System roles (r1, r2) tidak bisa dihapus/dimodifikasi

---

## 18. Data Flow Examples

### Attendance Flow
1. Karyawan check-in di fingerprint device
2. Device kirim NIK + timestamp ke backend
3. Backend validasi: NIK exists, tidak absent, FTW status fit
4. Backend determine status: `hadir` atau `terlambat`
5. Insert ke `attendance_logs`
6. Frontend poll `/api/attendance/today` untuk update UI

### Roster Upload Flow
1. Admin upload Excel roster via `/api/rosters/upload`
2. Backend parse Excel, validate rows
3. Create `roster_files` record
4. Create `roster_schedules` untuk setiap karyawan x tanggal
5. Seed `roster_revisions` jika ada perubahan dari bulan sebelumnya
6. Frontend display roster dengan color coding per shift

### Fleet Allocation Flow
1. Trigger: Auto allocate untuk tanggal dan shift tertentu
2. Backend cek unit dengan status `ready`
3. Assign unit ke fleet setting yang aktif
4. Create `fleet_allocations` record
5. Frontend display allocation dengan unit details
6. Update `unit_statuses` → `standby` atau `ready` sesuai allocation

---

## 19. Database Schema Highlights

### Core Tables
- `users`, `roles`, `role_permissions` → Auth & RBAC
- `employees`, `employee_competencies` → Karyawan
- `roster_files`, `roster_schedules`, `roster_revisions` → Roster
- `attendance_logs` → Absensi
- `ftw_logs` → Fit To Work
- `fleet_settings`, `fleet_setting_units`, `fleet_allocations` → Fleet
- `unit_statuses`, `unit_status_histories`, `units_db` → Assets
- `prestasi_scores`, `prestasi_badges`, `prestasi_history` → Prestasi
- `master_entries` → Master Data
- `notifications` → Notifications
- `app_settings`, `audio_schedules`, `audio_schedule_displays`, `display_devices` → Settings

---

## 20. API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Success message",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "errors": ["detail error 1", "detail error 2"]
}
```

### Status Codes
- `200`: Success
- `201`: Created
- `400`: Bad Request (validation error)
- `401`: Unauthorized (no/invalid token)
- `403`: Forbidden (insufficient permission)
- `404`: Not Found
- `500`: Internal Server Error

---

## 21. Security Considerations

1. **JWT Secret**: Diset via env `JWT_SECRET`, default development only
2. **Password Hashing**: Salt + hash,FE format untuk new users
3. **SQL Injection**: GORM parameterized queries
4. **CORS**: Configurable via `CORS_ALLOWED_ORIGINS`
5. **RBAC**: Server-side enforcement, frontend hanya UI hide
6. **File Upload**: Validasi file type dan size
7. **Environment Variables**: Sensitive config di `.env` (tidak di-commit)

---

## 22. Deployment & Operations

### Environment Variables Required
```
APP_NAME=universev2-backend
APP_ENV=development
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=universev2
JWT_SECRET=your-secret-key
FINGERPRINT_PORT=/dev/ttyUSB0
FINGERPRINT_BAUD=115200
FINGERPRINT_ENABLED=false
```

### Run Commands
```bash
# Development
make run

# With Docker
make dc-up

# Run migrations
make migrate-up

# Run tests
make test
```

---

## 23. Future Enhancements

1. **WebSocket**: Real-time attendance updates ke display devices
2. **Mobile App**: API untuk operator akses roster dan attendance
3. **Advanced Reporting**: Excel/PDF export dengan formatting
4. **Audit Log**: Track semua perubahan data penting
5. **Multi-company**: Support multiple perusahaan dalam satu instance
6. **Offline Mode**: Frontend cache untuk daerah dengan koneksi buruk

---

*Document generated based on backend code analysis and frontend module structure.*