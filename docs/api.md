# API Documentation

Dokumen ini merangkum endpoint utama setelah implementasi Fase 1–4 dan menjelaskan cara memakai dokumentasi Swagger yang sudah digenerate.

## Swagger Usage

### Interactive UI

Setelah aplikasi berjalan, buka:

```text
http://localhost:8080/swagger/index.html
```

Di Swagger UI Anda bisa:
- melihat request/response schema
- melihat field wajib, query params, dan auth requirement
- mencoba request langsung dari browser

### Regenerate Docs

Jika annotation di handler berubah, generate ulang file Swagger dengan:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs
go build ./cmd/api
```

Generated files:
- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

## Authentication

Semua endpoint protected menggunakan header:

```http
Authorization: Bearer <token>
```

---

## 1. Auth

### POST /login
Login user dan mengembalikan JWT token.

**Body**

```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

**Response**

```json
{
  "token": "<jwt>"
}
```

Gunakan token tersebut pada endpoint protected:

```http
Authorization: Bearer <token>
```

---

## 2. Locations

### GET /locations
Ambil semua lokasi, termasuk tipe lokasi, parameter geofence, dan minimum staffing.

### POST /locations
Buat lokasi baru untuk outlet atau warehouse.

**Body**

```json
{
  "name": "Warehouse Bekasi",
  "type": "warehouse",
  "address": "Jl. Example",
  "city": "Bekasi",
  "province": "Jawa Barat",
  "latitude": -6.2,
  "longitude": 106.9,
  "geofence_radius": 100,
  "minimum_staffing": 8
}
```

### GET /locations/{id}
Ambil detail satu lokasi.

### PUT /locations/{id}
Ubah data lokasi, termasuk `type`, koordinat, radius geofence, dan `minimum_staffing`.

### DELETE /locations/{id}
Hapus lokasi.

---

## 3. Users

### GET /users
Ambil daftar user beserta relasi lokasi/shift yang tersedia pada response.

### POST /users
Buat user baru.

**Body**

```json
{
  "name": "Budi",
  "email": "budi@example.com",
  "role": "staff_gudang",
  "location_id": 1,
  "shift_id": 2,
  "password": "secret",
  "base_salary": 5000000
}
```

### GET /users/{id}
Ambil detail satu user.

### PUT /users/{id}
Ubah role, lokasi, shift, password, atau gaji pokok user.

### DELETE /users/{id}
Hapus user jika tidak sedang memegang aset borrowed.

---

## 4. Attendance

### POST /attendance/checkin
Absensi masuk berbasis geofence, selfie, shift, dan compliance APD.

**Form Data**
- `location_id`
- `latitude`
- `longitude`
- `selfie` (file)

**Behaviour penting**
- lokasi harus sama dengan lokasi user
- jika lokasi bertipe `warehouse`, backend mengecek APD wajib
- jika shift user ada, status bisa `present` atau `late`

**Contoh error compliance**

```json
{
  "error": "check-in diblokir: APD wajib belum lengkap atau tidak layak"
}
```

### POST /attendance/checkout
Checkout attendance terbuka user. Jika checkout melewati akhir shift, payroll dapat menghitung `overtime_pay`.

### GET /attendance
Riwayat attendance user.

### GET /attendance/{id}
Detail attendance.

---

## 5. Leaves

### POST /leaves
Ajukan cuti / izin.

**Form Data**
- `start_date` (`YYYY-MM-DD`)
- `end_date` (`YYYY-MM-DD`)
- `leave_type`
- `reason`
- `document` (optional file)

### GET /leaves
Riwayat cuti user.

### GET /leaves/{id}
Detail cuti.

### POST /leaves/approve/{id}
Approve cuti dan kembalikan evaluasi staffing risk bila tersedia.

**Response**

```json
{
  "leave": {
    "id": 10,
    "status": "approved"
  },
  "staffing_risk": {
    "warning": true,
    "message": "Risiko Understaffed jika disetujui",
    "minimum_staffing": 8,
    "estimated_available_staff": 7
  }
}
```

### POST /leaves/reject/{id}
Reject cuti.

---

## 6. Payroll

### GET /payrolls?month={m}&year={y}
Hitung / ambil payroll user.

**Response fields penting**
- `base_salary`
- `cold_storage_allowance`
- `overtime_pay`
- `gross_salary`
- `total_deduction`
- `net_salary`

Jika `month` dan `year` tidak diisi, backend memakai periode bulan berjalan.

### GET /payrolls/history
Riwayat payroll user.

### GET /payrolls/download?month={m}&year={y}
Unduh slip gaji PDF.

---

## 7. Assets

### GET /assets
Ambil daftar aset/APD.

### POST /assets

**Body**

```json
{
  "name": "Jaket Thermal A1",
  "category": "jaket thermal",
  "serial_number": "APD-001",
  "condition": "Layak",
  "notes": "Untuk staff cold storage"
}
```

### GET /assets/{id}
Ambil detail satu aset.

### PUT /assets/{id}
Ubah field aset seperti `name`, `category`, `condition`, dan `notes`.

### DELETE /assets/{id}
Hapus aset.

### POST /assets/borrow/{id}
Assign aset ke seorang user.

**Body**

```json
{
  "user_id": 12
}
```

### POST /assets/return/{id}
Kembalikan aset dan set kondisi terakhirnya.

**Body**

```json
{
  "condition": "Layak"
}
```

---

## 8. Executive Reports

Endpoint report dibatasi untuk role `admin_hr`.

### GET /api/reports/labour-cost-leakage
Laporan eksekutif untuk melihat kebocoran biaya tenaga kerja akibat telat/mangkir per lokasi.

Query params optional:
- `month`
- `year`

**Response shape**

```json
{
  "data": [
    {
      "location_id": 1,
      "location_name": "Warehouse Bekasi",
      "location_type": "warehouse",
      "city": "Bekasi",
      "headcount": 20,
      "total_base_salary": 100000000,
      "total_gross_salary": 112000000,
      "total_deduction": 3500000,
      "total_net_salary": 108500000,
      "total_absent_count": 4,
      "total_late_count": 9,
      "leakage_percentage": 3.12
    }
  ]
}
```

### GET /api/reports/attendance-risk
Laporan eksekutif untuk melihat tekanan staffing dan anomali attendance per lokasi.

Query params optional:
- `month`
- `year`

**Response shape**

```json
{
  "data": [
    {
      "location_id": 1,
      "location_name": "Warehouse Bekasi",
      "location_type": "warehouse",
      "city": "Bekasi",
      "headcount": 20,
      "attendance_records": 350,
      "open_attendance_count": 2,
      "approved_leave_count": 3,
      "pending_leave_count": 1,
      "attendance_risk_score": 8.5,
      "staffing_risk": "medium",
      "minimum_staffing": 18,
      "estimated_available_staff": 17
    }
  ]
}
```

---

## 9. Swagger

### GET /swagger/index.html
Route UI Swagger yang sudah terdaftar pada aplikasi. Gunakan ini sebagai referensi interaktif utama untuk schema dan percobaan request.
