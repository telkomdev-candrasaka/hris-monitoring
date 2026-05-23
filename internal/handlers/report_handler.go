package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// LabourCostLeakageHandler godoc
// @Summary Get labour cost leakage report
// @Description Executive report showing attendance-related payroll leakage by location, including deductions, lateness, absence counts, and leakage percentage.
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param month query int false "Month"
// @Param year query int false "Year"
// @Success 200 {object} LabourCostLeakageReportListDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/reports/labour-cost-leakage [get]
func (h *ReportHandler) LabourCostLeakageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	month, year, err := h.parseMonthYear(r)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.service.GetLabourCostLeakage(month, year)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil labour cost leakage report"})
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{"data": result})
}

// AttendanceRiskHandler godoc
// @Summary Get attendance risk report
// @Description Executive report showing attendance anomalies, leave load, minimum staffing pressure, and computed attendance risk by location.
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param month query int false "Month"
// @Param year query int false "Year"
// @Success 200 {object} AttendanceRiskReportListDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/reports/attendance-risk [get]
func (h *ReportHandler) AttendanceRiskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	month, year, err := h.parseMonthYear(r)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.service.GetAttendanceRisk(month, year)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil attendance risk report"})
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{"data": result})
	}

func (h *ReportHandler) parseMonthYear(r *http.Request) (int, int, error) {
	month := 0
	year := 0
	query := r.URL.Query()
	if rawMonth := query.Get("month"); rawMonth != "" {
		parsed, err := strconv.Atoi(rawMonth)
		if err != nil || parsed < 1 || parsed > 12 {
			return 0, 0, errInvalid("month invalid")
		}
		month = parsed
	}
	if rawYear := query.Get("year"); rawYear != "" {
		parsed, err := strconv.Atoi(rawYear)
		if err != nil || parsed < 2000 {
			return 0, 0, errInvalid("year invalid")
		}
		year = parsed
	}
	return month, year, nil
}

type invalidError string

func errInvalid(message string) error {
	return invalidError(message)
}

func (e invalidError) Error() string {
	return string(e)
}

func (h *ReportHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
