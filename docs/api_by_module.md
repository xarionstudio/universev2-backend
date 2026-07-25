# UniverseV2 Backend API Documentation

Dokumentasi API backend berdasarkan modul-frontend.

---

## 1. Dashboard

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/dashboard/summary` | Ringkasan data dashboard |

---

## 2. Displays

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/settings/displays` | List display devices |
| POST | `/api/v1/settings/displays` | Create display device |
| PUT | `/api/v1/settings/displays/:id` | Update display device |
| DELETE | `/api/v1/settings/displays/:id` | Delete display device |
| GET | `/api/v1/displays/:id/heartbeat` | Get display heartbeat/status |

---

## 3. Employees

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/employees/` | List employees (filter by dept, status, search) |
| GET | `/api/v1/employees/:nik` | Get employee by NIK |
| POST | `/api/v1/employees/` | Create employee |
| PUT | `/api/v1/employees/:nik` | Update employee |
| DELETE | `/api/v1/employees/:nik` | Delete employee |
| GET | `/api/v1/employees/:nik/competencies` | Get employee competencies |
| PUT | `/api/v1/employees/:nik/competencies` | Update competencies |
| POST | `/api/v1/employees/:nik/photo` | Upload employee photo |

---

## 4. Fit To Work (FTW)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/ftw/today` | Get today's FTW logs |
| GET | `/api/v1/ftw/history` | Get FTW history (query: `nik`, `startDate`, `endDate`) |
| POST | `/api/v1/ftw/submit` | Submit FTW log |

---

## 5. Roster

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/rosters/` | List roster files |
| GET | `/api/v1/rosters/:key/export` | Export roster as Excel |
| GET | `/api/v1/rosters/:key/detail` | Get roster detail |
| POST | `/api/v1/rosters/upload` | Upload roster file |
| GET | `/api/v1/rosters/revisions` | Get revisions (query: `status`) |
| GET | `/api/v1/rosters/revisions/codes` | Get revision codes |
| POST | `/api/v1/rosters/revisions/batch` | Submit batch revision |
| DELETE | `/api/v1/rosters/revisions/:id` | Delete revision |
| PUT | `/api/v1/rosters/approvals/:id/approve` | Approve revision |
| PATCH | `/api/v1/rosters/approvals/:id/note` | Approve revision with note |
| PUT | `/api/v1/rosters/approvals/:id/reject` | Reject revision |
| GET | `/api/v1/rosters/attendance` | Get roster attendance (query: `date`, `dept`) |

---

## 6. Attendance

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/attendance/today` | Get today's attendance (query: `dept`, `shift`) |
| GET | `/api/v1/attendance/date` | Get attendance by date (query: `date`) |
| GET | `/api/v1/attendance/range` | Get attendance range (query: `startDate`, `endDate`) |
| POST | `/api/v1/attendance/checkin` | Record check-in |
| POST | `/api/v1/attendance/checkout` | Record check-out |

---

## 7. Fleet / Assets

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/fleet/settings` | Get fleet settings |
| POST | `/api/v1/fleet/settings` | Create fleet setting |
| PUT | `/api/v1/fleet/settings/:id` | Update fleet setting |
| DELETE | `/api/v1/fleet/settings/:id` | Delete fleet setting |
| GET | `/api/v1/fleet/allocations` | Get fleet allocations (query: `date`, `shift`) |
| POST | `/api/v1/fleet/allocations/auto` | Auto allocate fleet |
| GET | `/api/v1/units/status` | Get unit statuses |
| PUT | `/api/v1/units/:code/status` | Update unit status |
| POST | `/api/v1/units/:code/status-report` | Report unit breakdown |
| GET | `/api/v1/units/:code/history` | Get unit history |
| GET | `/api/v1/units/db` | Get unit DB |
| POST | `/api/v1/units/db` | Create unit DB |
| PUT | `/api/v1/units/db` | Update unit DB |
| DELETE | `/api/v1/units/db` | Delete unit DB (query: `id`) |

---

## 8. Prestasi (Performance)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/prestasi/leaderboard` | Get performance leaderboard (query: `periodDays`) |
| GET | `/api/v1/prestasi/:nik/history` | Get operator performance history (query: `periodDays`) |

---

## 9. Master Data

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/master/:category` | Get master entries by category |
| POST | `/api/v1/master/:category` | Create master entry |
| PUT | `/api/v1/master/:category/:id` | Update master entry |
| DELETE | `/api/v1/master/:category/:id` | Delete master entry |

---

## 10. Users & Roles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/` | List users |
| POST | `/api/v1/users/` | Create user |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |
| GET | `/api/v1/roles/` | List roles |
| POST | `/api/v1/roles/` | Create role |
| PUT | `/api/v1/roles/:id` | Update role |
| DELETE | `/api/v1/roles/:id` | Delete role |

---

## 11. Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/notifications/` | Get user notifications |
| PUT | `/api/v1/notifications/:id/read` | Mark notification as read |
| PUT | `/api/v1/notifications/read-all` | Mark all notifications as read |

---

## 12. Settings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/settings/` | Get app settings |
| PUT | `/api/v1/settings/` | Update app settings |
| GET | `/api/v1/settings/audio` | Get audio schedules |
| POST | `/api/v1/settings/audio` | Create audio schedule |
| PUT | `/api/v1/settings/audio/:id` | Update audio schedule |
| DELETE | `/api/v1/settings/audio/:id` | Delete audio schedule |
| GET | `/api/v1/settings/displays` | Get display devices |
| POST | `/api/v1/settings/displays` | Create display device |
| PUT | `/api/v1/settings/displays/:id` | Update display device |
| DELETE | `/api/v1/settings/displays/:id` | Delete display device |

---

## 13. Fingerprint

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/fingerprint/devices` | Get fingerprint devices status |

---

## 14. Auth (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/register` | Register |
| POST | `/api/v1/auth/refresh` | Refresh token |
| POST | `/api/v1/auth/logout` | Logout |

---

## 15. Profile

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/profile/` | Get profile |
| PUT | `/api/v1/profile/` | Update profile |
| PUT | `/api/v1/profile/password` | Change password |

---

## 16. Weather

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/weather/current` | Get current weather |

---

**Base URL:** `http://localhost:8080`  
**Auth:** Bearer JWT token required for most endpoints, except `/api/v1/auth/*`  
**RBAC:** Role-based access control applies per module (`view` / `manage`)