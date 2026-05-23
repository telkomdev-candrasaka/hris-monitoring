package services

import (
	"errors"
	"time"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"
)

type LeaveService struct {
	repo *repositories.LeaveRepository
	userRepo *repositories.UserRepository
}

func NewLeaveService(repo *repositories.LeaveRepository, userRepo *repositories.UserRepository) *LeaveService {
	return &LeaveService{repo: repo, userRepo: userRepo}
}

type StaffingRiskResult struct {
	Warning               bool   `json:"warning"`
	Message               string `json:"message,omitempty"`
	MinimumStaffing       int    `json:"minimum_staffing"`
	EstimatedAvailableStaff int64 `json:"estimated_available_staff"`
}

func (s *LeaveService) ApplyLeave(userID uint, startDate, endDate time.Time, leaveType, reason, documentPath string) (*models.Leave, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("tanggal akhir cuti tidak boleh sebelum tanggal mulai")
	}

	leave := &models.Leave{
		UserID:       userID,
		StartDate:    startDate,
		EndDate:      endDate,
		LeaveType:    leaveType,
		Reason:       reason,
		Status:       "pending",
		DocumentPath: documentPath,
	}

	if err := s.repo.CreateLeave(leave); err != nil {
		return nil, err
	}

	return leave, nil
}

func (s *LeaveService) GetLeaveHistory(userID uint) ([]models.Leave, error) {
	return s.repo.GetLeavesByUser(userID)
}

func (s *LeaveService) GetLeaveByID(id uint) (*models.Leave, error) {
	return s.repo.GetLeaveByID(id)
}

func (s *LeaveService) ApproveLeave(leaveID, approverID uint, action string) (*models.Leave, error) {
	leave, err := s.repo.GetLeaveByID(leaveID)
	if err != nil {
		return nil, err
	}

	if leave.Status != "pending" {
		return nil, errors.New("pengajuan cuti sudah diproses sebelumnya")
	}

	if _, err := s.EvaluateStaffingRisk(leave.UserID, leave.StartDate, leave.EndDate); err != nil {
		return nil, err
	}

	status := "approved"
	if action == "reject" {
		status = "rejected"
	} else if action != "approve" {
		return nil, errors.New("aksi tidak valid")
	}

	now := time.Now()
	leave.Status = status
	leave.ApprovedBy = approverID
	leave.ApprovedAt = &now

	if err := s.repo.UpdateLeave(leave); err != nil {
		return nil, err
	}

	return leave, nil
}

func (s *LeaveService) GetPendingLeaves() ([]models.Leave, error) {
	return s.repo.GetPendingLeaves()
}

func (s *LeaveService) EvaluateStaffingRisk(userID uint, startDate, endDate time.Time) (*StaffingRiskResult, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	headcount, err := s.userRepo.CountUsersByLocation(user.LocationID)
	if err != nil {
		return nil, err
	}

	approved, err := s.repo.CountApprovedLeavesByLocationAndDateRange(user.LocationID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	estimated := headcount - approved - 1
	if estimated < 0 {
		estimated = 0
	}

	result := &StaffingRiskResult{
		Warning: estimated < int64(user.Location.MinimumStaffing),
		MinimumStaffing: user.Location.MinimumStaffing,
		EstimatedAvailableStaff: estimated,
	}
	if result.Warning {
		result.Message = "Risiko Understaffed jika disetujui"
	}
	return result, nil
}
