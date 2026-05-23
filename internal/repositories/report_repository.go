package repositories

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
)

type LabourCostLeakageReportRow struct {
	LocationID        uint    `json:"location_id"`
	LocationName      string  `json:"location_name"`
	LocationType      string  `json:"location_type"`
	City              string  `json:"city"`
	Headcount         int64   `json:"headcount"`
	TotalBaseSalary   float64 `json:"total_base_salary"`
	TotalGrossSalary  float64 `json:"total_gross_salary"`
	TotalDeduction    float64 `json:"total_deduction"`
	TotalNetSalary    float64 `json:"total_net_salary"`
	TotalAbsentCount  int64   `json:"total_absent_count"`
	TotalLateCount    int64   `json:"total_late_count"`
	LeakagePercentage float64 `json:"leakage_percentage"`
}

type AttendanceRiskReportRow struct {
	LocationID                uint    `json:"location_id"`
	LocationName              string  `json:"location_name"`
	LocationType              string  `json:"location_type"`
	City                      string  `json:"city"`
	Headcount                 int64   `json:"headcount"`
	AttendanceRecords         int64   `json:"attendance_records"`
	OpenAttendanceCount       int64   `json:"open_attendance_count"`
	ApprovedLeaveCount        int64   `json:"approved_leave_count"`
	PendingLeaveCount         int64   `json:"pending_leave_count"`
	AttendanceRiskScore       float64 `json:"attendance_risk_score"`
	StaffingRisk              string  `json:"staffing_risk"`
	MinimumStaffing           int     `json:"minimum_staffing"`
	EstimatedAvailableStaff   int64   `json:"estimated_available_staff"`
}

type ReportRepository struct{}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{}
}

func (r *ReportRepository) GetLabourCostLeakage(month, year int) ([]LabourCostLeakageReportRow, error) {
	rows := []LabourCostLeakageReportRow{}
	query := config.DB.Table("payrolls p").
		Select(`
			l.id as location_id,
			l.name as location_name,
			l.type as location_type,
			l.city as city,
			count(distinct u.id) as headcount,
			coalesce(sum(p.base_salary), 0) as total_base_salary,
			coalesce(sum(p.gross_salary), 0) as total_gross_salary,
			coalesce(sum(p.total_deduction), 0) as total_deduction,
			coalesce(sum(p.net_salary), 0) as total_net_salary,
			coalesce(sum(p.absent_count), 0) as total_absent_count,
			coalesce(sum(p.late_count), 0) as total_late_count,
			case when coalesce(sum(p.gross_salary), 0) > 0 then (coalesce(sum(p.total_deduction), 0) / sum(p.gross_salary)) * 100 else 0 end as leakage_percentage`).
		Joins("join users u on u.id = p.user_id").
		Joins("join locations l on l.id = u.location_id")

	if month > 0 {
		query = query.Where("p.month = ?", month)
	}
	if year > 0 {
		query = query.Where("p.year = ?", year)
	}

	err := query.Group("l.id, l.name, l.type, l.city").Order("total_deduction desc").Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetAttendanceRisk(month, year int) ([]AttendanceRiskReportRow, error) {
	rows := []AttendanceRiskReportRow{}
	query := config.DB.Table("locations l").
		Select(`
			l.id as location_id,
			l.name as location_name,
			l.type as location_type,
			l.city as city,
			l.minimum_staffing as minimum_staffing,
			count(distinct u.id) as headcount,
			count(a.id) as attendance_records,
			coalesce(sum(case when a.check_out is null then 1 else 0 end), 0) as open_attendance_count,
			coalesce(sum(case when le.status = 'approved' then 1 else 0 end), 0) as approved_leave_count,
			coalesce(sum(case when le.status = 'pending' then 1 else 0 end), 0) as pending_leave_count,
			greatest(count(distinct u.id) - coalesce(sum(case when le.status = 'approved' then 1 else 0 end), 0), 0) as estimated_available_staff,
			(
				(coalesce(sum(case when a.check_out is null then 1 else 0 end), 0) * 2) +
				coalesce(sum(case when le.status = 'approved' then 1 else 0 end), 0) +
				(coalesce(sum(case when le.status = 'pending' then 1 else 0 end), 0) * 0.5)
			) as attendance_risk_score,
			case
				when greatest(count(distinct u.id) - coalesce(sum(case when le.status = 'approved' then 1 else 0 end), 0), 0) < l.minimum_staffing then 'high'
				when greatest(count(distinct u.id) - coalesce(sum(case when le.status = 'approved' then 1 else 0 end), 0), 0) = l.minimum_staffing then 'medium'
				else 'low'
			end as staffing_risk`).
		Joins("left join users u on u.location_id = l.id").
		Joins("left join attendances a on a.user_id = u.id").
		Joins("left join leaves le on le.user_id = u.id")

	if month > 0 {
		query = query.Where("a.id is null or extract(month from a.check_in) = ?", month)
		query = query.Where("le.id is null or extract(month from le.start_date) = ? or extract(month from le.end_date) = ?", month, month)
	}
	if year > 0 {
		query = query.Where("a.id is null or extract(year from a.check_in) = ?", year)
		query = query.Where("le.id is null or extract(year from le.start_date) = ? or extract(year from le.end_date) = ?", year, year)
	}

	err := query.Group("l.id, l.name, l.type, l.city, l.minimum_staffing").Order("attendance_risk_score desc").Scan(&rows).Error
	return rows, err
}
