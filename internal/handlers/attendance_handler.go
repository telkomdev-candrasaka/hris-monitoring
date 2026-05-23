package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/middlewares"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/services"
)

type AttendanceHandler struct {
	service *services.AttendanceService
}

type checkInResponse struct {
	Attendance interface{} `json:"attendance"`
}

type checkOutResponse struct {
	Attendance interface{} `json:"attendance"`
}

func NewAttendanceHandler(service *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

// CheckInHandler godoc
// @Summary Check in attendance
// @Description Check in with location validation, geofence validation, selfie upload, shift-aware lateness, and warehouse PPE compliance.
// @Tags Attendance
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param location_id formData int true "Location ID"
// @Param latitude formData number true "Device latitude"
// @Param longitude formData number true "Device longitude"
// @Param selfie formData file true "Selfie image"
// @Success 201 {object} AttendanceResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /attendance/checkin [post]
func (h *AttendanceHandler) CheckInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal memproses form data"})
		return
	}

	locationID, err := strconv.ParseUint(r.FormValue("location_id"), 10, 64)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "location_id tidak valid"})
		return
	}

	deviceLat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "latitude tidak valid"})
		return
	}

	deviceLng, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "longitude tidak valid"})
		return
	}

	file, header, err := r.FormFile("selfie")
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Selfie wajib diunggah"})
		return
	}
	defer file.Close()

	selfiePath, err := h.saveSelfieFile(file, header, claims.UserID)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menyimpan selfie"})
		return
	}

	attendance, err := h.service.CheckIn(claims.UserID, uint(locationID), deviceLat, deviceLng, selfiePath)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusCreated, attendance)
}

// CheckOutHandler godoc
// @Summary Check out attendance
// @Description Close the current open attendance record for the authenticated user.
// @Tags Attendance
// @Security BearerAuth
// @Produce json
// @Success 200 {object} AttendanceResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /attendance/checkout [post]
func (h *AttendanceHandler) CheckOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	attendance, err := h.service.CheckOut(claims.UserID)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusOK, attendance)
}

// GetAttendancesHandler godoc
// @Summary List current user's attendances
// @Description Return attendance history for the authenticated user.
// @Tags Attendance
// @Security BearerAuth
// @Produce json
// @Success 200 {array} AttendanceResponseDoc
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /attendance [get]
func (h *AttendanceHandler) GetAttendancesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	attendances, err := h.service.GetAttendancesByUser(claims.UserID)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data absensi"})
		return
	}

	h.writeJSON(w, http.StatusOK, attendances)
}

// GetAttendanceByIDHandler godoc
// @Summary Get attendance by ID
// @Description Retrieve one attendance record. Admin HR and outlet managers may access records beyond their own when authorized.
// @Tags Attendance
// @Security BearerAuth
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} AttendanceResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /attendance/{id} [get]
func (h *AttendanceHandler) GetAttendanceByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	id, err := h.parseID(r.URL.Path)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID absensi tidak valid"})
		return
	}

	attendance, err := h.service.GetAttendanceByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "Absensi tidak ditemukan"})
		return
	}

	if attendance.UserID != claims.UserID && claims.Role != "admin_hr" && claims.Role != "manager_outlet" {
		h.writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
		return
	}

	h.writeJSON(w, http.StatusOK, attendance)
}

func (h *AttendanceHandler) saveSelfieFile(file io.Reader, header *multipart.FileHeader, userID uint) (string, error) {
	destinationDir := "uploads/selfies"
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(destinationDir, filename)

	dest, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		return "", err
	}

	return filePath, nil
}

func (h *AttendanceHandler) parseID(path string) (uint, error) {
	trimmed := strings.TrimPrefix(path, "/attendance/")
	value, err := strconv.Atoi(trimmed)
	return uint(value), err
}

func (h *AttendanceHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
