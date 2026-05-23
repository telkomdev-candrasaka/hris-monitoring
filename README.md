# HRIS Monitoring & Operations API

Backend API berbasis Golang untuk HRIS dan operasional multi-cabang pada bisnis retail frozen food. Sistem ini mencakup autentikasi berbasis JWT, pengelolaan master data lokasi dan user, absensi geofencing dengan selfie, workflow cuti, payroll, pengelolaan aset/APD, kepatuhan K3 warehouse, serta dashboard eksekutif berbasis laporan agregasi.

## Overview

Proyek ini dibangun untuk membantu perusahaan retail frozen food mengelola disiplin kehadiran, staffing continuity, payroll berbasis kehadiran, serta kepatuhan APD pada lingkungan warehouse/cold storage. Implementasi saat ini menggunakan arsitektur berlapis `handler -> service -> repository -> PostgreSQL` dengan HTTP server berbasis `net/http` dan ORM GORM.

## Implemented Features

Seluruh poin berikut sudah terverifikasi dari kode yang ada saat ini.

### 1. Authentication & RBAC
- Login dengan email dan password pada `POST /login`
- Password di-hash menggunakan bcrypt
- JWT Bearer token berisi `user_id`, `role`, dan `email`
- Middleware autentikasi untuk endpoint terproteksi
- Pembatasan role untuk area administratif seperti user management, asset management, leave approval, dan executive reports

### 2. Location Management
- CRUD lokasi pada:
  - `GET /locations`
  - `POST /locations`
  - `GET /locations/{id}`
  - `PUT /locations/{id}`
  - `DELETE /locations/{id}`
- Lokasi sekarang menyimpan:
  - `type` (`outlet` / `warehouse`)
  - `minimum_staffing`
  - data geofence (`latitude`, `longitude`, `geofence_radius`)

### 3. User Management
- CRUD user pada:
  - `GET /users`
  - `POST /users`
  - `GET /users/{id}`
  - `PUT /users/{id}`
  - `DELETE /users/{id}`
- User dapat dikaitkan dengan:
  - lokasi kerja
  - shift kerja
  - gaji pokok (`base_salary`)
- Penghapusan user diblokir bila masih memiliki aset yang belum dikembalikan

### 4. Shift-Aware Attendance
- Check-in multipart dengan selfie pada `POST /attendance/checkin`
- Check-out pada `POST /attendance/checkout`
- Riwayat absensi user pada `GET /attendance`
- Detail absensi pada `GET /attendance/{id}`
- Absensi memvalidasi:
  - kecocokan lokasi kerja user
  - geofence lokasi
  - keterlambatan berdasarkan shift dan grace period bila shift user tersedia
- Status attendance kini mendukung alur seperti:
  - `present`
  - `late`
  - `present_checked_out`
  - `late_checked_out`
- Selfie disimpan ke `uploads/selfies`

### 5. Warehouse Safety Compliance (K3 / APD)
- Sistem mengenal `mandatory_equipments` per lokasi dan role
- Saat user warehouse melakukan check-in, backend memvalidasi apakah user memegang APD wajib dengan kondisi `Layak`
- Jika APD wajib tidak lengkap atau tidak layak, check-in diblokir
- Pengelolaan aset tetap tersedia untuk peminjaman/pengembalian APD dan aset operasional

### 6. Leave Workflow & Staffing Continuity
- Ajukan cuti/izin pada `POST /leaves`
- Riwayat cuti pribadi pada `GET /leaves`
- Detail cuti pada `GET /leaves/{id}`
- Approve cuti pada `POST /leaves/approve/{id}`
- Reject cuti pada `POST /leaves/reject/{id}`
- Upload dokumen pendukung opsional ke `uploads/leaves`
- Saat approval dilakukan, backend mengevaluasi staffing risk berdasarkan:
  - `minimum_staffing` lokasi
  - headcount lokasi
  - approved leave yang overlap
- Response approval dapat mengandung warning:
  - `Risiko Understaffed jika disetujui`

### 7. Payroll & Payslip
- Hitung / ambil payroll bulanan pada `GET /payrolls?month={m}&year={y}`
- Riwayat payroll pada `GET /payrolls/history`
- Download slip gaji PDF pada `GET /payrolls/download?month={m}&year={y}`
- Payroll saat ini menghitung:
  - `base_salary`
  - `cold_storage_allowance` untuk user di lokasi `warehouse`
  - `overtime_pay` jika checkout melebihi akhir shift
  - `gross_salary`
  - `total_deduction`
  - `net_salary`
- PDF payslip juga sudah menampilkan komponen payroll baru tersebut

### 8. Asset & PPE Management
- CRUD aset pada:
  - `GET /assets`
  - `POST /assets`
  - `GET /assets/{id}`
  - `PUT /assets/{id}`
  - `DELETE /assets/{id}`
- Peminjaman aset pada `POST /assets/borrow/{id}`
- Pengembalian aset pada `POST /assets/return/{id}`
- Aset memiliki status dan kondisi, termasuk default `Layak`
- Rule resign safeguard tetap berlaku: user tidak boleh dihapus bila masih memegang aset borrowed

### 9. Executive Reports
- `GET /api/reports/labour-cost-leakage`
- `GET /api/reports/attendance-risk`
- Report mengambil data agregat dari payroll, attendance, leave, user, dan location
- Mendukung filter `month` dan `year`
- Endpoint dibatasi untuk role `admin_hr`

### 10. Interactive Swagger Documentation
- Swagger UI tersedia di `/swagger/index.html`
- Endpoint sudah dilengkapi annotation untuk request/response utama
- File generated Swagger disimpan di:
  - `docs/docs.go`
  - `docs/swagger.json`
  - `docs/swagger.yaml`
- Dokumentasi interaktif dapat dipakai untuk:
  - melihat schema request/response
  - mencoba endpoint dari browser
  - memeriksa kebutuhan Bearer token dan query parameter

## Architecture

Sistem mengikuti layered architecture:

```text
Client
  -> Handler
    -> Service
      -> Repository
        -> PostgreSQL
```

### Business Rule Highlights
- `Location.Type` mengendalikan perbedaan rule outlet vs warehouse
- `MandatoryEquipment` + `Asset.Condition` mengendalikan kepatuhan APD warehouse saat check-in
- `Location.MinimumStaffing` dipakai untuk warning staffing continuity pada leave approval
- `Shift` dipakai untuk validasi keterlambatan dan perhitungan overtime payroll

## Project Structure

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── docs/
│   ├── docs.go
│   ├── architecture.md
│   ├── api.md
│   └── images/
│       ├── main.png
│       ├── feature.png
│       └── extra.png
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── middlewares/
│   ├── models/
│   ├── repositories/
│   └── services/
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Tech Stack

- Go 1.26.3
- `net/http`
- PostgreSQL
- GORM
- JWT (`github.com/golang-jwt/jwt/v5`)
- Bcrypt
- Godotenv
- gofpdf
- Swaggo HTTP Swagger route integration

## Data Model Summary

Entitas utama saat ini:

- `Location`
- `Shift`
- `User`
- `Attendance`
- `Leave`
- `Payroll`
- `Asset`
- `MandatoryEquipment`

Semua tabel dimigrasikan lewat `AutoMigrate` saat aplikasi start.

## Authentication

Gunakan Bearer token pada header:

```http
Authorization: Bearer <token>
```

Role yang tampak digunakan pada kode saat ini:
- `admin_hr`
- `manager_outlet`

## Environment Variables

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hris-monitoring
DB_PORT=5432
DB_SSLMODE=disable
JWT_SECRET=change-me-in-production
```

Catatan:
- `JWT_SECRET` wajib diisi untuk produksi
- aplikasi tetap fallback ke secret default jika env tidak tersedia, tetapi itu tidak disarankan untuk production

## Getting Started

### 1. Clone repository

```bash
git clone <repository-url>
cd hris
```

### 2. Siapkan environment

Buat `.env`:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hris-monitoring
DB_PORT=5432
DB_SSLMODE=disable
JWT_SECRET=change-me-in-production
```

### 3. Jalankan PostgreSQL

Pastikan database target tersedia dan kredensial sesuai.

### 4. Jalankan aplikasi

```bash
go run ./cmd/api
```

Atau build manual:

```bash
go build ./cmd/api
```

### 5. Akses service

- API: `http://localhost:8080`
- Swagger route: `http://localhost:8080/swagger/index.html`

### 6. Regenerate Swagger docs

Jika annotation handler berubah, generate ulang dokumentasi dengan:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs
```

Lalu verifikasi build:

```bash
go build ./cmd/api
```

## Running with Docker Compose

Repository ini menyediakan:
- `docker-compose.yml`
- `Dockerfile`

Jalankan:

```bash
docker compose up --build
```

Service utama:
- PostgreSQL di `5432`
- API di `8080`

## API Overview

### Public
- `POST /login`

### Protected Core Modules
- `/locations`
- `/users`
- `/attendance`
- `/leaves`
- `/payrolls`
- `/assets`

### Executive Reports
- `/api/reports/labour-cost-leakage`
- `/api/reports/attendance-risk`

Dokumentasi endpoint lebih rinci tersedia di `docs/api.md`.

Untuk dokumentasi interaktif, jalankan aplikasi lalu buka `http://localhost:8080/swagger/index.html`.

## Screenshots

### Main Interface
![Main](./docs/images/main.png)

### Feature View
![Feature](./docs/images/feature.png)

## Roadmap

Peningkatan realistis berikutnya:
- CRUD admin untuk `Shift` dan `MandatoryEquipment`
- contoh response yang lebih kaya untuk seluruh endpoint Swagger
- validasi request yang lebih ketat
- seed data awal untuk role, shift, dan mandatory equipments
- report dashboard yang lebih kaya per outlet vs warehouse
- audit trail untuk perubahan approval dan compliance

## Contributing

- Ikuti pola `handler -> service -> repository`
- Hindari mencampur business logic ke handler
- Jaga perubahan schema tetap additive bila memungkinkan
- Sinkronkan dokumentasi dengan perilaku aktual backend

## License

Belum ada file lisensi yang terdeteksi di repository ini.
