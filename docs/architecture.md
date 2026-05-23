# Architecture Overview

## 1. High-Level Architecture

Sistem menggunakan arsitektur backend berlapis untuk memisahkan concern HTTP, business rule, dan akses data.

```text
Client / Frontend
  -> HTTP Handlers
    -> Services
      -> Repositories
        -> PostgreSQL via GORM
```

Server dijalankan dari `cmd/api/main.go` menggunakan `net/http` dengan route registration manual.

## 2. Layers

### Handler Layer
Folder: `internal/handlers`

Tanggung jawab:
- menerima request HTTP
- membaca path/query/body/form-data
- mengembalikan response JSON atau file PDF

Handler yang tersedia saat ini mencakup:
- auth
- location
- user
- attendance
- leave
- payroll
- asset
- report

### Service Layer
Folder: `internal/services`

Tanggung jawab:
- business logic utama
- validasi operasional
- decision flow lintas entitas

Contoh rule penting:
- geofence attendance
- warehouse PPE compliance check
- staffing continuity warning pada leave approval
- shift-aware attendance status
- cold storage allowance dan overtime payroll

### Repository Layer
Folder: `internal/repositories`

Tanggung jawab:
- query database via GORM
- CRUD standar
- query agregasi untuk report
- query overlap leave dan lookup APD wajib

### Model Layer
Folder: `internal/models`

Model utama:
- `Location`
- `Shift`
- `User`
- `Attendance`
- `Leave`
- `Payroll`
- `Asset`
- `MandatoryEquipment`

## 3. Data Flow

### Authentication Flow
1. Client mengirim `POST /login`
2. Handler memanggil `AuthService`
3. Service memvalidasi password hash
4. Service menghasilkan JWT
5. Middleware `JWTAuth` memverifikasi token pada request berikutnya

### Attendance Flow
1. User kirim check-in multipart (`location_id`, `latitude`, `longitude`, `selfie`)
2. Handler membaca token, form, dan file
3. `AttendanceService` mengambil user dan location
4. Service memvalidasi kecocokan lokasi user
5. Jika lokasi bertipe `warehouse`, service memvalidasi APD wajib terhadap `MandatoryEquipment` dan `Asset.Condition`
6. Service memvalidasi geofence
7. Service menentukan status `present` atau `late` berdasarkan shift
8. Attendance disimpan ke database

### Leave Approval Flow
1. Manager/Admin membuka request leave
2. Approval endpoint memanggil `LeaveService.ApproveLeave`
3. Service mengecek status request
4. Service mengevaluasi staffing continuity berdasarkan `Location.MinimumStaffing`
5. Leave di-update menjadi `approved` atau `rejected`
6. Handler mengembalikan payload leave dan staffing risk warning jika relevan

### Payroll Flow
1. Payroll diminta berdasarkan user, month, year
2. `PayrollService` mengambil user dan attendance bulan terkait
3. Service menghitung deduction dari absent/late counters
4. Service menambahkan:
   - `cold_storage_allowance` untuk warehouse
   - `overtime_pay` dari checkout melebihi akhir shift
5. Payroll disimpan dan dapat diunduh sebagai PDF

### Report Flow
1. Admin HR memanggil endpoint report
2. `ReportService` memanggil repository agregasi
3. Repository menjalankan grouped query lintas payroll, attendance, leave, user, dan location
4. Handler mengembalikan data agregat siap dipakai dashboard

## 4. Realtime / Event Flow

Belum ada mekanisme realtime atau event bus pada implementasi saat ini.

Semua proses berjalan dalam pola request/response sinkron.

## 5. Key Design Decisions

### a. Location-Driven Rules
`Location.Type` digunakan sebagai pemicu rule bisnis untuk membedakan outlet dan warehouse.

### b. Safety Compliance as Check-In Gate
Warehouse compliance tidak hanya dicatat, tetapi dipakai sebagai guard sebelum check-in berhasil.

### c. Staffing Continuity as Approval Insight
Leave approval tidak di-hard-block, tetapi mengembalikan warning operasional agar keputusan tetap bisa diambil oleh manajer.

### d. Additive Schema Evolution
Perubahan fase 1–4 dilakukan secara additive agar kompatibel dengan pendekatan `AutoMigrate` yang digunakan repo ini.

### e. Report as Separate Read-Only Module
Executive reporting dipisahkan ke repository/service/handler tersendiri agar query agregat tidak mencemari CRUD repository biasa.
