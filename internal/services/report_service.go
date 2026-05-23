package services

import "github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetLabourCostLeakage(month, year int) ([]repositories.LabourCostLeakageReportRow, error) {
	return s.repo.GetLabourCostLeakage(month, year)
}

func (s *ReportService) GetAttendanceRisk(month, year int) ([]repositories.AttendanceRiskReportRow, error) {
	return s.repo.GetAttendanceRisk(month, year)
}
