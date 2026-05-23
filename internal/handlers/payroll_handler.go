package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/middlewares"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/services"
)

type PayrollHandler struct {
	service            *services.PayrollService
	pdfService         *services.PayslipPDFService
	userRepo           *repositories.UserRepository
}

func NewPayrollHandler(service *services.PayrollService) *PayrollHandler {
	return &PayrollHandler{service: service}
}

func NewPayrollHandlerWithPDF(service *services.PayrollService, pdfService *services.PayslipPDFService, userRepo *repositories.UserRepository) *PayrollHandler {
	return &PayrollHandler{service: service, pdfService: pdfService, userRepo: userRepo}
}

func (h *PayrollHandler) RegisterRoutes() {
	http.HandleFunc("/payrolls", h.PayrollHandler)
	http.HandleFunc("/payrolls/history", h.PayrollHistoryHandler)
}

// PayrollHandler godoc
// @Summary Calculate or get payroll
// @Description Calculate current user's payroll for a given month and year, including warehouse allowance and overtime.
// @Tags Payroll
// @Security BearerAuth
// @Produce json
// @Param month query int false "Month"
// @Param year query int false "Year"
// @Success 200 {object} PayrollResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /payrolls [get]
func (h *PayrollHandler) PayrollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	query := r.URL.Query()
	month := time.Now().Month()
	year := time.Now().Year()

	if m := query.Get("month"); m != "" {
		parsed, err := strconv.Atoi(m)
		if err != nil || parsed < 1 || parsed > 12 {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "month invalid"})
			return
		}
		month = time.Month(parsed)
	}

	if y := query.Get("year"); y != "" {
		parsed, err := strconv.Atoi(y)
		if err != nil || parsed < 2000 {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "year invalid"})
			return
		}
		year = parsed
	}

	payroll, err := h.service.CalculatePayroll(claims.UserID, int(month), year)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusOK, payroll)
}

// PayrollHistoryHandler godoc
// @Summary Get payroll history
// @Description Return previously generated payroll records for the authenticated user.
// @Tags Payroll
// @Security BearerAuth
// @Produce json
// @Success 200 {array} PayrollResponseDoc
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /payrolls/history [get]
func (h *PayrollHandler) PayrollHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	history, err := h.service.GetPayrollHistory(claims.UserID)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil riwayat payslip"})
		return
	}

	h.writeJSON(w, http.StatusOK, history)
}

// DownloadPayslipPDFHandler godoc
// @Summary Download payslip PDF
// @Description Generate and download a PDF payslip for the authenticated user and selected month/year.
// @Tags Payroll
// @Security BearerAuth
// @Produce application/pdf
// @Param month query int false "Month"
// @Param year query int false "Year"
// @Success 200 {file} file "Payslip PDF"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /payrolls/download [get]
func (h *PayrollHandler) DownloadPayslipPDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	if h.pdfService == nil || h.userRepo == nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "PDF service not available"})
		return
	}

	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	query := r.URL.Query()
	month := time.Now().Month()
	year := time.Now().Year()

	if m := query.Get("month"); m != "" {
		parsed, err := strconv.Atoi(m)
		if err != nil || parsed < 1 || parsed > 12 {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "month invalid"})
			return
		}
		month = time.Month(parsed)
	}

	if y := query.Get("year"); y != "" {
		parsed, err := strconv.Atoi(y)
		if err != nil || parsed < 2000 {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "year invalid"})
			return
		}
		year = parsed
	}

	payroll, err := h.service.CalculatePayroll(claims.UserID, int(month), year)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "User tidak ditemukan"})
		return
	}

	pdfBytes, err := h.pdfService.GeneratePayslipPDF(payroll, user)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal generate PDF"})
		return
	}

	filename := fmt.Sprintf("payslip_%d_%02d_%d.pdf", claims.UserID, month, year)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}

func (h *PayrollHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
