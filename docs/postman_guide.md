# Panduan Postman Collection - UniverseV2 Backend API

Dokumen ini berisi panduan lengkap cara menggunakan setiap endpoint API UniverseV2 di Postman.

## Daftar Isi

- [Konfigurasi Awal](#konfigurasi-awal)
- [1. Auth](#1-auth)
- [2. Profile](#2-profile)
- [3. Weather](#3-weather)
- [4. Employees](#4-employees)
- [5. Fit To Work (FTW)](#5-fit-to-work-ftw)
- [6. Roster](#6-roster)
- [7. Attendance](#7-attendance)
- [8. Fleet & Units](#8-fleet--units)
- [9. Prestasi](#9-prestasi)
- [10. Master Data](#10-master-data)
- [11. Users & Roles](#11-users--roles)
- [12. Notifications](#12-notifications)
- [13. Settings](#13-settings)
- [14. Fingerprint](#14-fingerprint)
- [Troubleshooting](#troubleshooting)

---

## Konfigurasi Awal

### Variabel Postman

Sebelum memulai, pastikan variable berikut sudah di-set di Postman:

| Variable | Nilai Default | Deskripsi |
|----------|--------------|-----------|
| `baseUrl` | `http://localhost:8080` | Base URL server API |
| `token` | `YOUR_JWT_TOKEN` | Token JWT untuk autentikasi |

### Cara Setting Variabel di Postman

1. Klik **"Environments"** di sidebar kiri Postman
2. Buat environment baru (misal: `UniverseV2 Local`)
3. Tambahkan variable:
   - `baseUrl` → Initial value: `http://localhost:8080`
   - `token` → Initial value: `eyJhbGciOi...` (token setelah login)
4. Pilih environment tersebut sebelum menjalankan request

### Autentikasi

Sebagian besar endpoint memerlukan **Bearer Token** (JWT). Ada 2 cara:

**Cara 1 - Collection-level Auth (Otomatis):**
Collection sudah di-set dengan auth type `Bearer Token` menggunakan `{{token}}`. Jika sudah mendapatkan token dan mengisi variable `token`, maka semua request yang membutuhkan auth akan otomatis menyertakan token.

**Cara 2 - Manual per Request:**
Tambahkan header:
```
Authorization: Bearer {{token}}
```

---

## 1. Auth

### 1.1 Login

> **Endpoint ini bersifat PUBLIC (tanpa token).**

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/auth/login` |
| **Auth** | `No Auth` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Cara menjalankan di Postman:**
1. Pilih method **POST**
2. Masukkan URL: `{{baseUrl}}/api/auth/login`
3. Di tab **Headers**, pastikan ada `Content-Type: application/json`
4. Di tab **Body**, pilih **raw** → **JSON**
5. Masukkan JSON body di atas
6. Klik **Send**

**Response sukses (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "u-xxxxx",
      "name": "John Doe",
      "email": "user@example.com",
      "nik": "503264133"
    }
  }
}
```

**Langkah setelah berhasil login:**
1. Copy value `token` dari response
2. Set variable `token` di Postman dengan token tersebut
3. Request selanjutnya akan otomatis menggunakan token

---

### 1.2 Register

> **Endpoint ini bersifat PUBLIC (tanpa token).**

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/auth/register` |
| **Auth** | `No Auth` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "name": "John Doe",
  "nik": "503264133",
  "email": "john@example.com",
  "password": "SecurePass123",
  "dept": "Operation",
  "pos": "Operator"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `name` | string | ✅ | Nama lengkap |
| `nik` | string | ✅ | NIK unik, validasi dari data employee |
| `email` | string | ✅ | Email unik untuk login |
| `password` | string | ✅ | Minimal 8 karakter |
| `dept` | string | ✅ | Departemen (harus sesuai data master) |
| `pos` | string | ✅ | Posisi/jabatan |

---

### 1.3 Refresh Token

> **Endpoint ini bersifat PUBLIC tapi membutuhkan header Authorization.**

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/auth/refresh` |
| **Auth** | `No Auth` (tapi perlu header Authorization manual) |
| **Content-Type** | `application/json` |

**Headers manual yang perlu ditambahkan:**
| Key | Value |
|-----|-------|
| `Content-Type` | `application/json` |
| `Authorization` | `Bearer {{token}}` |

**Request Body (raw JSON):**
```json
{}
```

> **Catatan:** Gunakan endpoint ini ketika token hampir expired. Response akan mengembalikan token baru.

---

### 1.4 Logout

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/auth/logout` |
| **Auth** | `Bearer Token` (otomatis via variable `{{token}}`) |
| **Content-Type** | `application/json` |

**Request Body:** Tidak perlu body.

---

## 2. Profile

### 2.1 Get Profile

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/profile/` |
| **Auth** | `Bearer Token` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/profile/`
3. Pastikan variable `{{token}}` sudah terisi
4. Klik **Send**

---

### 2.2 Update Profile

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/profile/` |
| **Auth** | `Bearer Token` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "name": "John Doe Updated"
}
```

**Field yang bisa di-update:**
| Field | Tipe | Keterangan |
|-------|------|------------|
| `name` | string | Nama lengkap |
| `email` | string | Email (jika diizinkan) |
| `nik` | string | NIK (jika diizinkan) |

---

### 2.3 Change Password

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/profile/password` |
| **Auth** | `Bearer Token` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "oldPassword": "oldpass123",
  "newPassword": "NewPass123",
  "confirmPassword": "NewPass123"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `oldPassword` | string | ✅ | Password saat ini |
| `newPassword` | string | ✅ | Password baru (min 8 karakter) |
| `confirmPassword` | string | ✅ | Harus sama dengan `newPassword` |

---

## 3. Weather

### 3.1 Get Current Weather

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/weather/current` |
| **Auth** | `Bearer Token` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/weather/current`
3. Klik **Send**

> **Catatan:** Endpoint ini mengambil data cuaca terkini dari API eksternal (BMKG).

---

## 4. Employees

### 4.1 List Employees

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/employees/` |
| **Auth** | `Bearer Token` |

**Query Parameters (opsional):**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `dept` | string | `Operation` | Filter berdasarkan departemen |
| `status` | string | `aktif` | Filter berdasarkan status (`aktif`/`nonaktif`) |
| `q` | string | `John` | Pencarian berdasarkan nama atau NIK |

**Contoh URL dengan parameter:**
```
{{baseUrl}}/api/employees/?dept=Operation&status=aktif&q=John
```

---

### 4.2 Get Employee by NIK

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/employees/:nik` |
| **Auth** | `Bearer Token` |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `nik` | string | `503264133` | NIK karyawan |

**Contoh URL:**
```
{{baseUrl}}/api/employees/503264133
```

---

### 4.3 Create Employee

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/employees/` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `employees:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "nik": "503264199",
  "name": "New Employee",
  "dept": "Operation",
  "pos": "Operator",
  "status": "aktif",
  "company": "PT Example"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `nik` | string | ✅ | NIK unik |
| `name` | string | ✅ | Nama lengkap |
| `dept` | string | ✅ | Departemen |
| `pos` | string | ✅ | Posisi/jabatan |
| `status` | string | ✅ | `aktif` / `nonaktif` |
| `company` | string | ❌ | Nama perusahaan (default: sesuai master) |

---

### 4.4 Update Employee

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/employees/:nik` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `employees:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `nik` | string | `503264133` |

**Request Body (raw JSON):**
```json
{
  "name": "Updated Name",
  "dept": "Plant",
  "pos": "Supervisor",
  "status": "aktif"
}
```

**Contoh URL:**
```
PUT {{baseUrl}}/api/employees/503264133
```

---

### 4.5 Delete Employee

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/employees/:nik` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `employees:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `nik` | string | `503264133` |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/employees/503264133
```

---

### 4.6 Get Competencies

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/employees/:nik/competencies` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `employees:view` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `nik` | string | `503264133` |

**Contoh URL:**
```
GET {{baseUrl}}/api/employees/503264133/competencies
```

---

### 4.7 Update Competencies

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/employees/:nik/competencies` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `employees:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `nik` | string | `503264133` |

**Request Body (raw JSON) — array of objects:**
```json
[
  {
    "className": "PC3000",
    "simperNo": "S-001",
    "expiryDate": "2026-12-31"
  }
]
```

**Deskripsi field per object:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `className` | string | ✅ | Nama kelas alat berat |
| `simperNo` | string | ✅ | Nomor SIMPER/SIO |
| `expiryDate` | string | ✅ | Tanggal kadaluarsa (format: `YYYY-MM-DD`) |

---

### 4.8 Upload Photo

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/employees/:nik/photo` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `employees:manage` |

**Cara upload file di Postman:**
1. Pilih method **POST**
2. Masukkan URL: `{{baseUrl}}/api/employees/503264133/photo`
3. Di tab **Body**, pilih **form-data**
4. Pada kolom **KEY**, isi `photo`
5. Ganti tipe dari **Text** menjadi **File** (dropdown di samping key)
6. Pada kolom **VALUE**, klik **Select Files** dan pilih file foto
7. Klik **Send**

**Deskripsi field form-data:**
| Key | Tipe | Value |
|-----|------|-------|
| `photo` | File | File gambar (jpg/png) |

---

## 5. Fit To Work (FTW)

### 5.1 Get Today FTW Log

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/ftw/today` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `ftw:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/ftw/today`
3. Klik **Send**

> **Catatan:** Endpoint ini mengembalikan log FTW untuk hari ini.

---

### 5.2 Get FTW History

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/ftw/history` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `ftw:view` |

**Query Parameters (opsional):**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `nik` | string | `503264133` | Filter berdasarkan NIK |
| `startDate` | string | `2026-07-01` | Filter tanggal awal (format: `YYYY-MM-DD`) |
| `endDate` | string | `2026-07-24` | Filter tanggal akhir (format: `YYYY-MM-DD`) |

**Contoh URL:**
```
{{baseUrl}}/api/ftw/history?nik=503264133&startDate=2026-07-01&endDate=2026-07-24
```

---

### 5.3 Submit FTW Log

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/ftw/submit` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `ftw:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "nik": "503264133",
  "shift": "siang",
  "sleepMinutes": 420,
  "sleepFormatted": "7 j 0 m",
  "restHours": 0,
  "canWork": true,
  "sendTime": "05:30"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `nik` | string | ✅ | NIK karyawan |
| `shift` | string | ✅ | Shift kerja (`pagi`/`siang`/`malam`) |
| `sleepMinutes` | number | ✅ | Durasi tidur dalam menit |
| `sleepFormatted` | string | ✅ | Durasi tidur dalam format `X j Y m` |
| `restHours` | number | ✅ | Jam istirahat (0 jika tidak ada) |
| `canWork` | boolean | ✅ | `true` jika fit bekerja, `false` jika tidak |
| `sendTime` | string | ✅ | Waktu pengiriman (format: `HH:mm`) |

---

## 6. Roster

### 6.1 List Rosters

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/rosters/` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/rosters/`
3. Klik **Send**

---

### 6.2 Export Roster

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/rosters/:key/export` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `key` | string | `jul` | Periode roster (contoh: `jul` untuk Juli) |

**Contoh URL:**
```
{{baseUrl}}/api/rosters/jul/export
```

> **Catatan:** Response berupa file Excel (.xlsx).

---

### 6.3 Get Roster Detail

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/rosters/:key/detail` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `key` | string | `jul` | Periode roster |

**Contoh URL:**
```
{{baseUrl}}/api/rosters/jul/detail
```

---

### 6.4 Upload Roster

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/rosters/upload` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:manage` |

**Cara upload file di Postman:**
1. Pilih method **POST**
2. Masukkan URL: `{{baseUrl}}/api/rosters/upload`
3. Di tab **Body**, pilih **form-data**
4. Pada kolom **KEY**, isi `file`
5. Ganti tipe dari **Text** menjadi **File**
6. Pada kolom **VALUE**, klik **Select Files** dan pilih file Excel roster
7. Klik **Send**

**Deskripsi field form-data:**
| Key | Tipe | Value |
|-----|------|-------|
| `file` | File | File Excel (.xlsx) berisi data roster |

---

### 6.5 Get Revisions

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/rosters/revisions` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Query Parameters (opsional):**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `status` | string | `pending` | Filter status (`pending`/`approved`/`rejected`) |

**Contoh URL:**
```
{{baseUrl}}/api/rosters/revisions?status=pending
```

---

### 6.6 Get Revision Codes

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/rosters/revisions/codes` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/rosters/revisions/codes`
3. Klik **Send**

> **Catatan:** Endpoint ini mengembalikan daftar kode yang digunakan untuk revisi roster (alasan perubahan).

---

### 6.7 Submit Batch Revision

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/rosters/revisions/batch` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON) — array of objects:**
```json
[
  {
    "submissionId": "sub-101",
    "nik": "503264133",
    "whatId": "Ganti shift D ke N",
    "whatEn": "Shift change D to N",
    "whenId": "Tgl 15 Jul",
    "whenEn": "15 Jul"
  }
]
```

**Deskripsi field per object:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `submissionId` | string | ✅ | ID unik pengajuan revisi |
| `nik` | string | ✅ | NIK karyawan |
| `whatId` | string | ✅ | Deskripsi perubahan (Bahasa Indonesia) |
| `whatEn` | string | ✅ | Deskripsi perubahan (English) |
| `whenId` | string | ✅ | Waktu perubahan (Bahasa Indonesia) |
| `whenEn` | string | ✅ | Waktu perubahan (English) |

---

### 6.8 Delete Revision

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/rosters/revisions/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `id` | string | `sub-101` | Submission ID revisi |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/rosters/revisions/sub-101
```

---

### 6.9 Approve Revision

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/rosters/approvals/:id/approve` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `sub-102` |

**Contoh URL:**
```
PUT {{baseUrl}}/api/rosters/approvals/sub-102/approve
```

**Cara menjalankan:**
1. Pilih method **PUT**
2. Masukkan URL: `{{baseUrl}}/api/rosters/approvals/sub-102/approve`
3. Klik **Send**

---

### 6.10 Approve Revision With Note

| Item | Detail |
|------|--------|
| **Method** | `PATCH` |
| **URL** | `{{baseUrl}}/api/rosters/approvals/:id/note` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `sub-102` |

**Request Body (raw JSON):**
```json
{
  "note": "Disetujui dengan catatan"
}
```

**Contoh URL:**
```
PATCH {{baseUrl}}/api/rosters/approvals/sub-102/note
```

---

### 6.11 Reject Revision

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/rosters/approvals/:id/reject` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `sub-101` |

**Contoh URL:**
```
PUT {{baseUrl}}/api/rosters/approvals/sub-101/reject
```

---

### 6.12 Get Roster Attendance

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/rosters/attendance` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Query Parameters:**
| Parameter | Tipe | Wajib | Contoh | Deskripsi |
|-----------|------|-------|--------|-----------|
| `date` | string | ❌ | `2026-07-24` | Tanggal (format: `YYYY-MM-DD`) |
| `dept` | string | ❌ | `Operation` | Filter departemen |

**Contoh URL:**
```
{{baseUrl}}/api/rosters/attendance?date=2026-07-24&dept=Operation
```

---

## 7. Attendance

### 7.1 Get Today Attendance

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/attendance/today` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Query Parameters (opsional):**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `dept` | string | `Operation` | Filter departemen |
| `shift` | string | `D` | Filter shift |

**Contoh URL:**
```
{{baseUrl}}/api/attendance/today?dept=Operation&shift=D
```

---

### 7.2 Get Attendance By Date

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/attendance/date` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Query Parameters:**
| Parameter | Tipe | Wajib | Contoh |
|-----------|------|-------|--------|
| `date` | string | ✅ | `2026-07-24` |

**Contoh URL:**
```
{{baseUrl}}/api/attendance/date?date=2026-07-24
```

---

### 7.3 Get Attendance Range

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/attendance/range` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:view` |

**Query Parameters:**
| Parameter | Tipe | Wajib | Contoh |
|-----------|------|-------|--------|
| `startDate` | string | ✅ | `2026-07-01` |
| `endDate` | string | ✅ | `2026-07-24` |

**Contoh URL:**
```
{{baseUrl}}/api/attendance/range?startDate=2026-07-01&endDate=2026-07-24
```

---

### 7.4 Record Check-In

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/attendance/checkin` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "nik": "503264133",
  "time": "05:45",
  "machine": "FP-01"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `nik` | string | ✅ | NIK karyawan |
| `time` | string | ✅ | Waktu check-in (format: `HH:mm`) |
| `machine` | string | ✅ | Kode mesin fingerprint |

---

### 7.5 Record Check-Out

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/attendance/checkout` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `roster:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "nik": "503264133",
  "time": "17:30",
  "machine": "FP-01"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `nik` | string | ✅ | NIK karyawan |
| `time` | string | ✅ | Waktu check-out (format: `HH:mm`) |
| `machine` | string | ✅ | Kode mesin fingerprint |

---

## 8. Fleet & Units

### 8.1 Get Fleet Settings

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/fleet/settings` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/fleet/settings`
3. Klik **Send**

---

### 8.2 Create Fleet Setting

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/fleet/settings` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "digger": "EX9001",
  "loc": "Workshop New",
  "bus": "BUS-01",
  "units": ["DT9001", "DT9002"]
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `digger` | string | ✅ | Kode excavator/digger |
| `loc` | string | ✅ | Lokasi |
| `bus` | string | ✅ | Kode bus (unit pendukung) |
| `units` | array[string] | ✅ | Daftar unit (dump truck, dll) |

---

### 8.3 Update Fleet Setting

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/fleet/settings/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `id` | string | `fl-EX7001` | ID fleet setting |

**Request Body (raw JSON):**
```json
{
  "digger": "EX7001",
  "loc": "Workshop Plant",
  "bus": "UDBU002",
  "units": ["DT114"]
}
```

**Contoh URL:**
```
PUT {{baseUrl}}/api/fleet/settings/fl-EX7001
```

---

### 8.4 Delete Fleet Setting

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/fleet/settings/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `fl-EX7001` |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/fleet/settings/fl-EX7001
```

---

### 8.5 Get Allocations

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/fleet/allocations` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:view` |

**Query Parameters (opsional):**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `date` | string | `2026-07-24` | Tanggal (format: `YYYY-MM-DD`) |
| `shift` | string | `pagi` | Shift (`pagi`/`siang`/`malam`) |

**Contoh URL:**
```
{{baseUrl}}/api/fleet/allocations?date=2026-07-24&shift=pagi
```

---

### 8.6 Auto Allocate

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/fleet/allocations/auto` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "date": "2026-07-24",
  "shift": "pagi"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `date` | string | ✅ | Tanggal (format: `YYYY-MM-DD`) |
| `shift` | string | ✅ | Shift (`pagi`/`siang`/`malam`) |

---

### 8.7 Get Unit Statuses

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/units/status` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/units/status`
3. Klik **Send**

---

### 8.8 Update Unit Status

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/units/:code/status` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `code` | string | `DT114` | Kode unit |

**Request Body (raw JSON):**
```json
{
  "status": "ready",
  "note": "Checklist harian ok"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `status` | string | ✅ | Status unit (`ready`/`breakdown`/`standby`/dll) |
| `note` | string | ❌ | Catatan tambahan |

**Contoh URL:**
```
PUT {{baseUrl}}/api/units/DT114/status
```

---

### 8.9 Report Unit Breakdown

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/units/:code/status-report` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `code` | string | `DT114` |

**Request Body (raw JSON):**
```json
{
  "reason": "Hydraulic sistem bocor"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `reason` | string | ✅ | Alasan/alasan breakdown |

**Contoh URL:**
```
POST {{baseUrl}}/api/units/DT114/status-report
```

---

### 8.10 Get Unit History

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/units/:code/history` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:view` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `code` | string | `DT114` |

**Contoh URL:**
```
{{baseUrl}}/api/units/DT114/history
```

---

### 8.11 Get Unit DB

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/units/db` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/units/db`
3. Klik **Send**

> **Catatan:** Endpoint ini mengembalikan daftar semua unit yang terdaftar di database (master data unit).

---

### 8.12 Create Unit DB

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/units/db` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "code": "DT9001",
  "egi": "EX900",
  "product": "Caterpillar",
  "className": "PC4000"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `code` | string | ✅ | Kode unit unik |
| `egi` | string | ❌ | Kode EGI/group |
| `product` | string | ❌ | Merek/produk |
| `className` | string | ❌ | Nama kelas alat |

---

### 8.13 Update Unit DB

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/units/db` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "code": "DT9001",
  "egi": "EX900",
  "product": "Caterpillar",
  "className": "PC4000",
  "category": "dt",
  "location": "Workshop"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `code` | string | ✅ | Kode unit (digunakan sebagai identifier) |
| `egi` | string | ❌ | Kode EGI/group |
| `product` | string | ❌ | Merek/produk |
| `className` | string | ❌ | Nama kelas alat |
| `category` | string | ❌ | Kategori (`dt` untuk dump truck, `ex` untuk excavator, dll) |
| `location` | string | ❌ | Lokasi unit |

---

### 8.14 Delete Unit DB

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/units/db?id=:code` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `asset:manage` |

**Query Parameter:**
| Parameter | Tipe | Wajib | Contoh |
|-----------|------|-------|--------|
| `id` | string | ✅ | `DT9001` |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/units/db?id=DT9001
```

---

## 9. Prestasi

### 9.1 Get Leaderboard

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/prestasi/leaderboard` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `prestasi:view` |

**Query Parameters (opsional):**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `periodDays` | number | `30` | Periode hari yang ditampilkan (default: 30) |

**Contoh URL:**
```
{{baseUrl}}/api/prestasi/leaderboard?periodDays=30
```

---

### 9.2 Get Operator History

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/prestasi/:nik/history` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `prestasi:view` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `nik` | string | `503264133` |

**Query Parameters (opsional):**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `periodDays` | number | `30` | Periode hari (default: 30) |

**Contoh URL:**
```
{{baseUrl}}/api/prestasi/503264133/history?periodDays=30
```

---

## 10. Master Data

Master data bersifat dinamis dengan kategori yang bisa disesuaikan (dept, pos, company, dll).

### 10.1 Get Master By Category

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/master/:category` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `master:view` |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `category` | string | `dept` | Kategori master (`dept`, `pos`, `company`, dll) |

**Contoh URL:**
```
{{baseUrl}}/api/master/dept
```

---

### 10.2 Create Master Entry

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/master/:category` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `master:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `category` | string | `dept` |

**Request Body (raw JSON):**
```json
{
  "name": "New Department",
  "fieldA": "",
  "fieldB": ""
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `name` | string | ✅ | Nama entry master |
| `fieldA` | string | ❌ | Field tambahan (tergantung kategori) |
| `fieldB` | string | ❌ | Field tambahan (tergantung kategori) |

---

### 10.3 Update Master Entry

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/master/:category/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `master:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `category` | string | `dept` |
| `id` | string | `dept-1` |

**Request Body (raw JSON):**
```json
{
  "name": "Department Updated",
  "fieldA": "",
  "fieldB": ""
}
```

**Contoh URL:**
```
PUT {{baseUrl}}/api/master/dept/dept-1
```

---

### 10.4 Delete Master Entry

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/master/:category/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `master:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `category` | string | `dept` |
| `id` | string | `dept-1` |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/master/dept/dept-1
```

---

## 11. Users & Roles

### 11.1 List Users

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/users/` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `users:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/users/`
3. Klik **Send**

---

### 11.2 Create User

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/users/` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `users:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "name": "New User",
  "email": "newuser@example.com",
  "password": "SecurePass123",
  "nik": "503264199",
  "roles": ["r1"]
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `name` | string | ✅ | Nama lengkap |
| `email` | string | ✅ | Email unik |
| `password` | string | ✅ | Password (min 8 karakter) |
| `nik` | string | ✅ | NIK (harus dari data employee yang sudah ada) |
| `roles` | array[string] | ❌ | Daftar role ID (contoh: `["r1", "r2"]`) |

---

### 11.3 Update User

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/users/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `users:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `u-1234567890` |

**Request Body (raw JSON):**
```json
{
  "name": "Updated User",
  "email": "updated@example.com",
  "roles": ["r1", "r2"]
}
```

**Contoh URL:**
```
PUT {{baseUrl}}/api/users/u-1234567890
```

---

### 11.4 Delete User

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/users/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `users:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `u-1234567890` |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/users/u-1234567890
```

---

### 11.5 Get Roles

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/roles/` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `users:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/roles/`
3. Klik **Send**

---

### 11.6 Create Role

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/roles/` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `users:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "id": "r6",
  "name": "Custom Role",
  "description": "Custom operator role"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `id` | string | ✅ | ID role unik |
| `name` | string | ✅ | Nama role |
| `description` | string | ❌ | Deskripsi role |

---

### 11.7 Update Role

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/roles/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `users:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `r6` |

**Request Body (raw JSON):**
```json
{
  "name": "Custom Role Updated",
  "description": "Updated description"
}
```

**Contoh URL:**
```
PUT {{baseUrl}}/api/roles/r6
```

---

### 11.8 Delete Role

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/roles/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `users:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `r6` |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/roles/r6
```

---

## 12. Notifications

### 12.1 Get Notifications

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/notifications/` |
| **Auth** | `Bearer Token` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/notifications/`
3. Klik **Send**

---

### 12.2 Mark Notification Read

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/notifications/:id/read` |
| **Auth** | `Bearer Token` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `n1` |

**Contoh URL:**
```
PUT {{baseUrl}}/api/notifications/n1/read
```

---

### 12.3 Mark All Read

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/notifications/read-all` |
| **Auth** | `Bearer Token` |

**Contoh URL:**
```
PUT {{baseUrl}}/api/notifications/read-all
```

---

## 13. Settings

### 13.1 Get Settings

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/settings/` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/settings/`
3. Klik **Send**

---

### 13.2 Update Settings

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/settings/` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "appName": "UNIVERSE-2",
  "theme": "dark",
  "lang": "id"
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `appName` | string | ✅ | Nama aplikasi |
| `theme` | string | ❌ | Tema (`light`/`dark`) |
| `lang` | string | ❌ | Bahasa (`id`/`en`) |

---

### 13.3 Get Audio Schedules

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/settings/audio` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/settings/audio`
3. Klik **Send**

---

### 13.4 Create Audio Schedule

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/settings/audio` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "title": "Pengumuman",
  "triggerTime": "05:45",
  "frequency": "harian",
  "fileName": "announcement.mp3",
  "isActive": true
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `title` | string | ✅ | Judul jadwal audio |
| `triggerTime` | string | ✅ | Waktu pemutaran (format: `HH:mm`) |
| `frequency` | string | ✅ | Frekuensi (`harian`/`mingguan`/dll) |
| `fileName` | string | ✅ | Nama file audio (harus ada di storage) |
| `isActive` | boolean | ✅ | Status aktif/nonaktif |

---

### 13.5 Update Audio Schedule

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/settings/audio/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `au1` |

**Request Body (raw JSON):**
```json
{
  "title": "Pengumuman P5M",
  "triggerTime": "05:45",
  "frequency": "harian",
  "fileName": "p5m_reminder.mp3",
  "isActive": true
}
```

**Contoh URL:**
```
PUT {{baseUrl}}/api/settings/audio/au1
```

---

### 13.6 Delete Audio Schedule

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/settings/audio/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `au1` |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/settings/audio/au1
```

---

### 13.7 Get Display Devices

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/settings/displays` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:view` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/settings/displays`
3. Klik **Send**

---

### 13.8 Create Display

| Item | Detail |
|------|--------|
| **Method** | `POST` |
| **URL** | `{{baseUrl}}/api/settings/displays` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:manage` |
| **Content-Type** | `application/json` |

**Request Body (raw JSON):**
```json
{
  "name": "TV Baru",
  "location": "Gate 2",
  "contentKind": "att",
  "runningText": "Patuhi batas kecepatan",
  "isActive": true
}
```

**Deskripsi field:**
| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `name` | string | ✅ | Nama display/TV |
| `location` | string | ✅ | Lokasi pemasangan |
| `contentKind` | string | ✅ | Jenis konten (`att` untuk attendance, `ftw`, dll) |
| `runningText` | string | ❌ | Teks berjalan (running text) |
| `isActive` | boolean | ✅ | Status aktif/nonaktif |

---

### 13.9 Update Display

| Item | Detail |
|------|--------|
| **Method** | `PUT` |
| **URL** | `{{baseUrl}}/api/settings/displays/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:manage` |
| **Content-Type** | `application/json` |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `id` | string | `DSP-A01` | ID display device |

**Request Body (raw JSON):**
```json
{
  "name": "TV Gate Utara",
  "location": "Gate utara",
  "contentKind": "att",
  "runningText": "Utamakan keselamatan",
  "isActive": true
}
```

**Contoh URL:**
```
PUT {{baseUrl}}/api/settings/displays/DSP-A01
```

---

### 13.10 Delete Display

| Item | Detail |
|------|--------|
| **Method** | `DELETE` |
| **URL** | `{{baseUrl}}/api/settings/displays/:id` |
| **Auth** | `Bearer Token` |
| **RBAC** | Membutuhkan permission `settings:manage` |

**Path Parameter:**
| Parameter | Tipe | Contoh |
|-----------|------|--------|
| `id` | string | `DSP-A01` |

**Contoh URL:**
```
DELETE {{baseUrl}}/api/settings/displays/DSP-A01
```

---

### 13.11 Get Display Heartbeat

> **Endpoint ini bersifat PUBLIC (tanpa token).**

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/displays/:id/heartbeat` |
| **Auth** | `No Auth` (public endpoint) |

**Path Parameter:**
| Parameter | Tipe | Contoh | Deskripsi |
|-----------|------|--------|-----------|
| `id` | string | `DSP-A01` | ID display device |

**Contoh URL:**
```
{{baseUrl}}/api/displays/DSP-A01/heartbeat
```

> **Catatan:** Endpoint ini digunakan oleh perangkat display TV untuk mengirim heartbeat (status online). Tidak memerlukan autentikasi.

---

## 14. Fingerprint

### 14.1 Get Fingerprint Devices

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **URL** | `{{baseUrl}}/api/fingerprint/devices` |
| **Auth** | `Bearer Token` |

**Cara menjalankan:**
1. Pilih method **GET**
2. Masukkan URL: `{{baseUrl}}/api/fingerprint/devices`
3. Klik **Send**

> **Catatan:** Endpoint ini mengembalikan daftar dan status mesin fingerprint yang terdaftar.

---

## Troubleshooting

### Masalah: Response 401 Unauthorized

**Penyebab:** Token tidak valid atau expired.
**Solusi:**
1. Login ulang menggunakan endpoint **Login**
2. Copy token baru
3. Update variable `{{token}}` di Postman
4. Coba request ulang

### Masalah: Response 403 Forbidden

**Penyebab:** User tidak memiliki permission yang diperlukan.
**Solusi:**
1. Cek RBAC permission yang tertera di tabel endpoint
2. Minta admin untuk memberikan role dengan permission yang sesuai

### Masalah: Response 400 Bad Request

**Penyebab:** Request body tidak valid atau field wajib tidak diisi.
**Solusi:**
1. Periksa kembali JSON body yang dikirim
2. Pastikan format tanggal sesuai (`YYYY-MM-DD`)
3. Pastikan semua field wajib terisi

### Masalah: Response 404 Not Found

**Penyebab:** Resource dengan ID/NIK yang diberikan tidak ditemukan.
**Solusi:**
1. Periksa kembali path parameter (NIK, ID, dll)
2. Pastikan data sudah ada di database

### Masalah: Response 500 Internal Server Error

**Penyebab:** Error di sisi server.
**Solusi:**
1. Hubungi developer
2. Periksa log server untuk detail error

### Tips Umum

1. **Gunakan folder Postman** - Kumpulan API sudah di-group per modul untuk memudahkan navigasi
2. **Simpan environment** - Gunakan environment berbeda untuk development, staging, dan production
3. **Test flow lengkap** - Urutan yang disarankan:
   - Login → copy token
   - Coba GET endpoints (profil, employees, dll)
   - Coba POST untuk membuat data baru
   - Coba PUT/PATCH untuk mengupdate
   - Coba DELETE untuk menghapus
4. **File upload** - Untuk endpoint upload file, jangan lupa ganti tipe dari **Text** ke **File** di form-data