package handlers

import "github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginRequestDoc struct {
	Email    string `json:"email" example:"admin@example.com"`
	Password string `json:"password" example:"secret123"`
}

type LoginResponseDoc struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

type LocationPayloadDoc struct {
	Name            string  `json:"name" example:"Warehouse Bekasi"`
	Type            string  `json:"type" example:"warehouse"`
	Address         string  `json:"address" example:"Jl. Raya Bekasi No. 1"`
	City            string  `json:"city" example:"Bekasi"`
	Province        string  `json:"province" example:"Jawa Barat"`
	Latitude        float64 `json:"latitude" example:"-6.2"`
	Longitude       float64 `json:"longitude" example:"106.9"`
	GeofenceRadius  float64 `json:"geofence_radius" example:"100"`
	MinimumStaffing int     `json:"minimum_staffing" example:"8"`
}

type LocationResponseDoc struct {
	ID              uint    `json:"id" example:"1"`
	Name            string  `json:"name" example:"Warehouse Bekasi"`
	Type            string  `json:"type" example:"warehouse"`
	Address         string  `json:"address" example:"Jl. Raya Bekasi No. 1"`
	City            string  `json:"city" example:"Bekasi"`
	Province        string  `json:"province" example:"Jawa Barat"`
	Latitude        float64 `json:"latitude" example:"-6.2"`
	Longitude       float64 `json:"longitude" example:"106.9"`
	GeofenceRadius  float64 `json:"geofence_radius" example:"100"`
	MinimumStaffing int     `json:"minimum_staffing" example:"8"`
}

type CreateUserRequestDoc struct {
	Name       string  `json:"name" example:"Budi"`
	Email      string  `json:"email" example:"budi@example.com"`
	Role       string  `json:"role" example:"staff_gudang"`
	LocationID uint    `json:"location_id" example:"1"`
	ShiftID    *uint   `json:"shift_id,omitempty" example:"2"`
	Password   string  `json:"password" example:"secret123"`
	BaseSalary float64 `json:"base_salary" example:"5000000"`
}

type UpdateUserRequestDoc struct {
	Name       string  `json:"name,omitempty" example:"Budi Update"`
	Email      string  `json:"email,omitempty" example:"budi.updated@example.com"`
	Role       string  `json:"role,omitempty" example:"manager_outlet"`
	LocationID uint    `json:"location_id,omitempty" example:"1"`
	ShiftID    *uint   `json:"shift_id,omitempty" example:"2"`
	Password   string  `json:"password,omitempty" example:"newsecret123"`
	BaseSalary float64 `json:"base_salary,omitempty" example:"6500000"`
}

type ShiftResponseDoc struct {
	ID            uint   `json:"id" example:"2"`
	LocationID    uint   `json:"location_id" example:"1"`
	Name          string `json:"name" example:"Shift Pagi"`
	StartTime     string `json:"start_time" example:"08:00"`
	EndTime       string `json:"end_time" example:"17:00"`
	CrossMidnight bool   `json:"cross_midnight" example:"false"`
	GraceMinutes  int    `json:"grace_minutes" example:"15"`
}

type UserResponseDoc struct {
	ID         uint                 `json:"id" example:"10"`
	Name       string               `json:"name" example:"Budi"`
	Email      string               `json:"email" example:"budi@example.com"`
	Role       string               `json:"role" example:"staff_gudang"`
	LocationID uint                 `json:"location_id" example:"1"`
	ShiftID    *uint                `json:"shift_id,omitempty" example:"2"`
	BaseSalary float64              `json:"base_salary" example:"5000000"`
	Location   *LocationResponseDoc `json:"location,omitempty"`
	Shift      *ShiftResponseDoc    `json:"shift,omitempty"`
}

type AttendanceResponseDoc struct {
	ID              uint    `json:"id" example:"1"`
	UserID          uint    `json:"user_id" example:"10"`
	LocationID      uint    `json:"location_id" example:"1"`
	Status          string  `json:"status" example:"present_checked_out"`
	CheckIn         string  `json:"check_in" example:"2026-05-23T08:05:00Z"`
	CheckOut        string  `json:"check_out,omitempty" example:"2026-05-23T17:30:00Z"`
	DeviceLatitude  float64 `json:"device_latitude" example:"-6.2"`
	DeviceLongitude float64 `json:"device_longitude" example:"106.9"`
	SelfiePath      string  `json:"selfie_path,omitempty" example:"uploads/selfies/user_10_123.jpg"`
}

type LeaveResponseDoc struct {
	ID           uint   `json:"id" example:"5"`
	UserID       uint   `json:"user_id" example:"10"`
	StartDate    string `json:"start_date" example:"2026-06-01T00:00:00Z"`
	EndDate      string `json:"end_date" example:"2026-06-03T00:00:00Z"`
	LeaveType    string `json:"leave_type" example:"sakit"`
	Reason       string `json:"reason,omitempty" example:"Demam tinggi"`
	Status       string `json:"status" example:"approved"`
	DocumentPath string `json:"document_path,omitempty" example:"uploads/leaves/leave_10_123.pdf"`
	ApprovedBy   uint   `json:"approved_by,omitempty" example:"1"`
	ApprovedAt   string `json:"approved_at,omitempty" example:"2026-05-23T10:00:00Z"`
}

type StaffingRiskResponseDoc struct {
	Warning                 bool   `json:"warning" example:"true"`
	Message                 string `json:"message,omitempty" example:"Risiko Understaffed jika disetujui"`
	MinimumStaffing         int    `json:"minimum_staffing" example:"8"`
	EstimatedAvailableStaff int64  `json:"estimated_available_staff" example:"7"`
}

type LeaveApprovalResponseDoc struct {
	Leave        LeaveResponseDoc          `json:"leave"`
	StaffingRisk *StaffingRiskResponseDoc  `json:"staffing_risk,omitempty"`
}

type CreateAssetRequestDoc struct {
	Name         string `json:"name" example:"Jaket Thermal A1"`
	Category     string `json:"category" example:"jaket thermal"`
	SerialNumber string `json:"serial_number" example:"APD-001"`
	Condition    string `json:"condition,omitempty" example:"Layak"`
	Notes        string `json:"notes,omitempty" example:"Untuk staff cold storage"`
}

type BorrowAssetRequestDoc struct {
	UserID uint `json:"user_id" example:"12"`
}

type ReturnAssetRequestDoc struct {
	Condition string `json:"condition" example:"Layak"`
}

type UpdateAssetRequestDoc struct {
	Name      string `json:"name,omitempty" example:"Jaket Thermal A2"`
	Category  string `json:"category,omitempty" example:"jaket thermal"`
	Condition string `json:"condition,omitempty" example:"Layak"`
	Notes     string `json:"notes,omitempty" example:"Dipakai untuk shift malam"`
}

type AssetResponseDoc struct {
	ID           uint   `json:"id" example:"1"`
	Name         string `json:"name" example:"Jaket Thermal A1"`
	Category     string `json:"category" example:"jaket thermal"`
	SerialNumber string `json:"serial_number" example:"APD-001"`
	UserID       *uint  `json:"user_id,omitempty" example:"10"`
	Status       string `json:"status" example:"borrowed"`
	Condition    string `json:"condition,omitempty" example:"Layak"`
	Notes        string `json:"notes,omitempty" example:"Untuk staff cold storage"`
}

type PayrollResponseDoc struct {
	ID                   uint    `json:"id" example:"1"`
	UserID               uint    `json:"user_id" example:"10"`
	Month                int     `json:"month" example:"5"`
	Year                 int     `json:"year" example:"2026"`
	BaseSalary           float64 `json:"base_salary" example:"5000000"`
	ColdStorageAllowance float64 `json:"cold_storage_allowance" example:"500000"`
	OvertimePay          float64 `json:"overtime_pay" example:"125000"`
	GrossSalary          float64 `json:"gross_salary" example:"5625000"`
	TotalAttendance      int     `json:"total_attendance" example:"22"`
	AbsentCount          int     `json:"absent_count" example:"0"`
	LateCount            int     `json:"late_count" example:"1"`
	TotalDeduction       float64 `json:"total_deduction" example:"50000"`
	NetSalary            float64 `json:"net_salary" example:"5575000"`
	GeneratedAt          string  `json:"generated_at" example:"2026-05-23T10:00:00Z"`
}

type ReportListResponseDoc struct {
	Data interface{} `json:"data"`
}

type LabourCostLeakageReportListDoc struct {
	Data []repositories.LabourCostLeakageReportRow `json:"data"`
}

type AttendanceRiskReportListDoc struct {
	Data []repositories.AttendanceRiskReportRow `json:"data"`
}
