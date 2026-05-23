package services

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"
)

type AttendanceService struct {
	repo         *repositories.AttendanceRepository
	locationRepo *repositories.LocationRepository
	userRepo     *repositories.UserRepository
	assetRepo    *repositories.AssetRepository
	mandatoryEquipmentRepo *repositories.MandatoryEquipmentRepository
}

func NewAttendanceService(repo *repositories.AttendanceRepository, locationRepo *repositories.LocationRepository, userRepo *repositories.UserRepository, assetRepo *repositories.AssetRepository, mandatoryEquipmentRepo *repositories.MandatoryEquipmentRepository) *AttendanceService {
	return &AttendanceService{repo: repo, locationRepo: locationRepo, userRepo: userRepo, assetRepo: assetRepo, mandatoryEquipmentRepo: mandatoryEquipmentRepo}
}

func (s *AttendanceService) haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000.0 // meters
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func (s *AttendanceService) CheckIn(userID, locationID uint, deviceLat, deviceLng float64, selfiePath string) (*models.Attendance, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	location, err := s.locationRepo.GetLocationByID(locationID)
	if err != nil {
		return nil, err
	}

	if user.LocationID != locationID {
		return nil, errors.New("lokasi absensi tidak sesuai dengan lokasi kerja user")
	}

	if err := s.validateWarehouseCompliance(user); err != nil {
		return nil, err
	}

	status, err := s.determineAttendanceStatus(user)
	if err != nil {
		return nil, err
	}

	distance := s.haversine(deviceLat, deviceLng, location.Latitude, location.Longitude)
	if distance > location.GeofenceRadius {
		return nil, errors.New("absensi ditolak: berada di luar radius geofence")
	}

	attendance := &models.Attendance{
		UserID:         userID,
		LocationID:     locationID,
		Status:         status,
		CheckIn:        time.Now(),
		DeviceLatitude:  deviceLat,
		DeviceLongitude: deviceLng,
		SelfiePath:     selfiePath,
	}

	if err := s.repo.CreateAttendance(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *AttendanceService) GetOpenAttendance(userID uint) (*models.Attendance, error) {
	return s.repo.GetOpenAttendance(userID)
}

func (s *AttendanceService) CheckOut(userID uint) (*models.Attendance, error) {
	attendance, err := s.repo.GetOpenAttendance(userID)
	if err != nil {
		return nil, err
	}

	if attendance.CheckOut != nil {
		return nil, errors.New("absensi sudah selesai")
	}

	now := time.Now()
	attendance.CheckOut = &now
	if attendance.Status == "late" {
		attendance.Status = "late_checked_out"
	} else {
		attendance.Status = "present_checked_out"
	}
	if err := s.repo.UpdateAttendance(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *AttendanceService) GetAttendancesByUser(userID uint) ([]models.Attendance, error) {
	return s.repo.GetAttendancesByUser(userID)
}

func (s *AttendanceService) GetAttendanceByID(id uint) (*models.Attendance, error) {
	return s.repo.GetAttendanceByID(id)
}

func (s *AttendanceService) validateWarehouseCompliance(user *models.User) error {
	if user.Location.Type != "warehouse" {
		return nil
	}

	requirements, err := s.mandatoryEquipmentRepo.GetByLocationAndRole(user.LocationID, user.Role)
	if err != nil {
		return err
	}
	if len(requirements) == 0 {
		return nil
	}

	assets, err := s.assetRepo.GetBorrowedAssetsByUser(user.ID)
	if err != nil {
		return err
	}

	owned := map[string]bool{}
	for _, asset := range assets {
		if strings.EqualFold(asset.Condition, "Layak") {
			owned[strings.ToLower(asset.Category)] = true
		}
	}

	for _, requirement := range requirements {
		if !owned[strings.ToLower(requirement.EquipmentCategory)] {
			return errors.New("check-in diblokir: APD wajib belum lengkap atau tidak layak")
		}
	}

	return nil
}

func (s *AttendanceService) determineAttendanceStatus(user *models.User) (string, error) {
	if user.Shift == nil {
		return "present", nil
	}

	now := time.Now()
	shiftStart, err := time.Parse("15:04", user.Shift.StartTime)
	if err != nil {
		return "present", err
	}
	shiftDateTime := time.Date(now.Year(), now.Month(), now.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location())
	graceDeadline := shiftDateTime.Add(time.Duration(user.Shift.GraceMinutes) * time.Minute)
	if now.After(graceDeadline) {
		return "late", nil
	}
	return "present", nil
}
