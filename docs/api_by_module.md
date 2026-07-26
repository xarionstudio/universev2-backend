# UniverseV2 Backend API Documentation

Dokumentasi API backend berdasarkan modul-frontend.

---

## 1. Dashboard

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/dashboard/summary` | Ringkasan data dashboard |

---

## 2. Displays

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/settings/displays` | List display devices |
| POST | `/api/settings/displays` | Create display device |
| PUT | `/api/settings/displays/:id` | Update display device |
| DELETE | `/api/settings/displays/:id` | Delete display device |
| GET | `/api/displays/:id/heartbeat` | Get display heartbeat/status |

---

## 3. Employees

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/employees/` | List employees (filter by dept, status, search) |
| GET | `/api/employees/:nik` | Get employee by NIK |
| POST | `/api/employees/` | Create employee |
| PUT | `/api/employees/:nik` | Update employee |
| DELETE | `/api/employees/:nik` | Delete employee |
| GET | `/api/employees/:nik/competencies` | Get employee competencies |
| PUT | `/api/employees/:nik/competencies` | Update competencies |
| POST | `/api/employees/:nik/photo` | Upload employee photo |

---

## 4. Fit To Work (FTW)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/ftw/today` | Get today's FTW logs |
| GET | `/api/ftw/history` | Get FTW history (query: `nik`, `startDate`, `endDate`) |
| POST | `/api/ftw/submit` | Submit FTW log |

---

## 5. Roster

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/rosters/` | List roster files |
| GET | `/api/rosters/:key/export` | Export roster as Excel |
| GET | `/api/rosters/:key/detail` | Get roster detail |
| POST | `/api/rosters/upload` | Upload roster file |
| GET | `/api/rosters/revisions` | Get revisions (query: `status`) |
| GET | `/api/rosters/revisions/codes` | Get revision codes |
| POST | `/api/rosters/revisions/batch` | Submit batch revision |
| DELETE | `/api/rosters/revisions/:id` | Delete revision |
| PUT | `/api/rosters/approvals/:id/approve` | Approve revision |
| PATCH | `/api/rosters/approvals/:id/note` | Approve revision with note |
| PUT | `/api/rosters/approvals/:id/reject` | Reject revision |
| GET | `/api/rosters/attendance` | Get roster attendance (query: `date`, `dept`) |

---

## 6. Attendance

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/attendance/today` | Get today's attendance (query: `dept`, `shift`) |
| GET | `/api/attendance/date` | Get attendance by date (query: `date`) |
| GET | `/api/attendance/range` | Get attendance range (query: `startDate`, `endDate`) |
| POST | `/api/attendance/checkin` | Record check-in |
| POST | `/api/attendance/checkout` | Record check-out |

---

## 7. Fleet / Assets

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/fleet/settings` | Get fleet settings |
| POST | `/api/fleet/settings` | Create fleet setting |
| PUT | `/api/fleet/settings/:id` | Update fleet setting |
| DELETE | `/api/fleet/settings/:id` | Delete fleet setting |
| GET | `/api/fleet/allocations` | Get fleet allocations (query: `date`, `shift`) |
| POST | `/api/fleet/allocations/auto` | Auto allocate fleet |
| GET | `/api/units/status` | Get unit statuses |
| PUT | `/api/units/:code/status` | Update unit status |
| POST | `/api/units/:code/status-report` | Report unit breakdown |
| GET | `/api/units/:code/history` | Get unit history |
| GET | `/api/units/db` | Get unit DB |
| POST | `/api/units/db` | Create unit DB |
| PUT | `/api/units/db` | Update unit DB |
| DELETE | `/api/units/db` | Delete unit DB (query: `id`) |

---

## 8. Prestasi (Performance)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/prestasi/leaderboard` | Get performance leaderboard (query: `periodDays`) |
| GET | `/api/prestasi/:nik/history` | Get operator performance history (query: `periodDays`) |

---

## 9. Master Data

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/master/:category` | Get master entries by category |
| POST | `/api/master/:category` | Create master entry |
| PUT | `/api/master/:category/:id` | Update master entry |
| DELETE | `/api/master/:category/:id` | Delete master entry |

---

## 10. Users & Roles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users/` | List users |
| POST | `/api/users/` | Create user |
| PUT | `/api/users/:id` | Update user |
| DELETE | `/api/users/:id` | Delete user |
| GET | `/api/roles/` | List roles |
| POST | `/api/roles/` | Create role |
| PUT | `/api/roles/:id` | Update role |
| DELETE | `/api/roles/:id` | Delete role |

---

## 11. Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/notifications/` | Get user notifications |
| PUT | `/api/notifications/:id/read` | Mark notification as read |
| PUT | `/api/notifications/read-all` | Mark all notifications as read |

---

## 12. Settings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/settings/` | Get app settings |
| PUT | `/api/settings/` | Update app settings |
| GET | `/api/settings/audio` | Get audio schedules |
| POST | `/api/settings/audio` | Create audio schedule |
| PUT | `/api/settings/audio/:id` | Update audio schedule |
| DELETE | `/api/settings/audio/:id` | Delete audio schedule |
| GET | `/api/settings/displays` | Get display devices |
| POST | `/api/settings/displays` | Create display device |
| PUT | `/api/settings/displays/:id` | Update display device |
| DELETE | `/api/settings/displays/:id` | Delete display device |

---

## 13. Fingerprint

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/fingerprint/devices` | Get fingerprint devices status |

---

## 14. Auth (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/register` | Register |
| POST | `/api/auth/refresh` | Refresh token |
| POST | `/api/auth/logout` | Logout |

---

## 15. Profile

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/profile/` | Get profile |
| PUT | `/api/profile/` | Update profile |
| PUT | `/api/profile/password` | Change password |

---

## 16. Weather

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/weather/current` | Get current weather |

---

**Base URL:** `http://localhost:8080`  
**Auth:** Bearer JWT token required for most endpoints, except `/api/auth/*`  
**RBAC:** Role-based access control applies per module (`view` / `manage`)