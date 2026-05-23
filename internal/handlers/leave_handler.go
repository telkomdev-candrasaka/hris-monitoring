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

type LeaveHandler struct {
	service *services.LeaveService
}

type leaveApprovalResponse struct {
	Leave        interface{}                    `json:"leave"`
	StaffingRisk *services.StaffingRiskResult   `json:"staffing_risk,omitempty"`
}

type applyLeaveResponse struct {
	Leave interface{} `json:"leave"`
}

func NewLeaveHandler(service *services.LeaveService) *LeaveHandler {
	return &LeaveHandler{service: service}
}

// applyLeave godoc
// @Summary Apply leave
// @Description Submit leave or permit request with optional supporting document upload.
// @Tags Leaves
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param start_date formData string true "Start date (YYYY-MM-DD)"
// @Param end_date formData string true "End date (YYYY-MM-DD)"
// @Param leave_type formData string true "Leave type"
// @Param reason formData string false "Reason"
// @Param document formData file false "Supporting document"
// @Success 201 {object} LeaveResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaves [post]
func (h *LeaveHandler) LeavesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.applyLeave(w, r)
	case http.MethodGet:
		h.getLeaveHistory(w, r)
	default:
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

// LeaveByIDHandler godoc
// @Summary Get leave by ID
// @Description Retrieve one leave request by identifier.
// @Tags Leaves
// @Security BearerAuth
// @Produce json
// @Param id path int true "Leave ID"
// @Success 200 {object} LeaveResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /leaves/{id} [get]
func (h *LeaveHandler) LeaveByIDHandler(w http.ResponseWriter, r *http.Request) {
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
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID cuti tidak valid"})
		return
	}

	leave, err := h.service.GetLeaveByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "Pengajuan cuti tidak ditemukan"})
		return
	}

	if leave.UserID != claims.UserID && claims.Role != "admin_hr" && claims.Role != "manager_outlet" {
		h.writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
		return
	}

	h.writeJSON(w, http.StatusOK, leave)
}

func (h *LeaveHandler) applyLeave(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal memproses form data"})
		return
	}

	startDate, err := time.Parse("2006-01-02", r.FormValue("start_date"))
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "start_date tidak valid, gunakan format YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", r.FormValue("end_date"))
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "end_date tidak valid, gunakan format YYYY-MM-DD"})
		return
	}

	leaveType := r.FormValue("leave_type")
	if leaveType == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "leave_type wajib diisi"})
		return
	}

	reason := r.FormValue("reason")
	documentPath := ""

	file, header, err := r.FormFile("document")
	if err == nil {
		defer file.Close()
		documentPath, err = h.saveDocumentFile(file, header, claims.UserID)
		if err != nil {
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menyimpan dokumen"})
			return
		}
	} else if err != http.ErrMissingFile {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal membaca file dokumen"})
		return
	}

	leave, err := h.service.ApplyLeave(claims.UserID, startDate, endDate, leaveType, reason, documentPath)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusCreated, leave)
}

// getLeaveHistory godoc
// @Summary List current user's leaves
// @Description Return leave history for the authenticated user.
// @Tags Leaves
// @Security BearerAuth
// @Produce json
// @Success 200 {array} LeaveResponseDoc
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leaves [get]
func (h *LeaveHandler) getLeaveHistory(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	leaves, err := h.service.GetLeaveHistory(claims.UserID)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil riwayat cuti"})
		return
	}

	h.writeJSON(w, http.StatusOK, leaves)
}

// ApproveLeaveHandler godoc
// @Summary Approve leave
// @Description Approve a pending leave request and return staffing risk evaluation.
// @Tags Leaves
// @Security BearerAuth
// @Produce json
// @Param id path int true "Leave ID"
// @Success 200 {object} LeaveApprovalResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /leaves/approve/{id} [post]
func (h *LeaveHandler) ApproveLeaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	id, err := h.parseApprovalID(r.URL.Path, "/leaves/approve/")
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID cuti tidak valid"})
		return
	}

	leave, err := h.service.ApproveLeave(id, claims.UserID, "approve")
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	staffingRisk, _ := h.service.EvaluateStaffingRisk(leave.UserID, leave.StartDate, leave.EndDate)
	h.writeJSON(w, http.StatusOK, leaveApprovalResponse{Leave: leave, StaffingRisk: staffingRisk})
}

// RejectLeaveHandler godoc
// @Summary Reject leave
// @Description Reject a pending leave request without staffing-risk payload.
// @Tags Leaves
// @Security BearerAuth
// @Produce json
// @Param id path int true "Leave ID"
// @Success 200 {object} LeaveApprovalResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /leaves/reject/{id} [post]
func (h *LeaveHandler) RejectLeaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	id, err := h.parseApprovalID(r.URL.Path, "/leaves/reject/")
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID cuti tidak valid"})
		return
	}

	leave, err := h.service.ApproveLeave(id, claims.UserID, "reject")
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusOK, leaveApprovalResponse{Leave: leave})
}

func (h *LeaveHandler) saveDocumentFile(file io.Reader, header *multipart.FileHeader, userID uint) (string, error) {
	destinationDir := "uploads/leaves"
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("leave_%d_%d%s", userID, time.Now().UnixNano(), ext)
	documentPath := filepath.Join(destinationDir, filename)

	dest, err := os.Create(documentPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		return "", err
	}

	return documentPath, nil
}

func (h *LeaveHandler) parseID(path string) (uint, error) {
	trimmed := strings.TrimPrefix(path, "/leaves/")
	value, err := strconv.Atoi(trimmed)
	return uint(value), err
}

func (h *LeaveHandler) parseApprovalID(path, prefix string) (uint, error) {
	trimmed := strings.TrimPrefix(path, prefix)
	value, err := strconv.Atoi(trimmed)
	return uint(value), err
}

func (h *LeaveHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
