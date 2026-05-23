package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/services"
)

type LocationHandler struct {
	service *services.LocationService
}

type locationPayload struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Address         string  `json:"address"`
	City            string  `json:"city"`
	Province        string  `json:"province"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	GeofenceRadius  float64 `json:"geofence_radius"`
	MinimumStaffing int     `json:"minimum_staffing"`
}

func NewLocationHandler(service *services.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (h *LocationHandler) RegisterRoutes() {
	http.HandleFunc("/locations", h.locationsHandler)
	http.HandleFunc("/locations/", h.locationByIDHandler)
}

func (h *LocationHandler) LocationsHandler(w http.ResponseWriter, r *http.Request) {
	h.locationsHandler(w, r)
}

func (h *LocationHandler) LocationByIDHandler(w http.ResponseWriter, r *http.Request) {
	h.locationByIDHandler(w, r)
}

// getAllLocations godoc
// @Summary List locations
// @Description Get all registered locations including type, geofence, and minimum staffing.
// @Tags Locations
// @Security BearerAuth
// @Produce json
// @Success 200 {array} LocationResponseDoc
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /locations [get]
func (h *LocationHandler) locationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllLocations(w, r)
	case http.MethodPost:
		h.createLocation(w, r)
	default:
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (h *LocationHandler) locationByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r.URL.Path)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID lokasi tidak valid"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getLocationByID(w, r, id)
	case http.MethodPut:
		h.updateLocation(w, r, id)
	case http.MethodDelete:
		h.deleteLocation(w, r, id)
	default:
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (h *LocationHandler) parseID(path string) (uint, error) {
	trimmed := strings.TrimPrefix(path, "/locations/")
	value, err := strconv.Atoi(trimmed)
	return uint(value), err
}

// createLocation godoc
// @Summary Create location
// @Description Create a new outlet or warehouse location used for geofencing, staffing rules, and location-driven payroll logic.
// @Tags Locations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body LocationPayloadDoc true "Location payload"
// @Success 201 {object} LocationResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /locations [post]
func (h *LocationHandler) createLocation(w http.ResponseWriter, r *http.Request) {
	var payload locationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}

	location := models.Location{
		Name: payload.Name,
		Type: payload.Type,
		Address: payload.Address,
		City: payload.City,
		Province: payload.Province,
		Latitude: payload.Latitude,
		Longitude: payload.Longitude,
		GeofenceRadius: payload.GeofenceRadius,
		MinimumStaffing: payload.MinimumStaffing,
	}
	if location.Type == "" {
		location.Type = "outlet"
	}

	if err := h.service.CreateLocation(&location); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal membuat lokasi"})
		return
	}

	h.writeJSON(w, http.StatusCreated, location)
}

func (h *LocationHandler) getAllLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.service.GetAllLocations()
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil daftar lokasi"})
		return
	}

	h.writeJSON(w, http.StatusOK, locations)
}

// getLocationByID godoc
// @Summary Get location by ID
// @Description Retrieve one location record including geofence and staffing fields.
// @Tags Locations
// @Security BearerAuth
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} LocationResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /locations/{id} [get]
func (h *LocationHandler) getLocationByID(w http.ResponseWriter, r *http.Request, id uint) {
	location, err := h.service.GetLocationByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "Lokasi tidak ditemukan"})
		return
	}

	h.writeJSON(w, http.StatusOK, location)
}

// updateLocation godoc
// @Summary Update location
// @Description Update an existing location including type, coordinates, geofence radius, and minimum staffing.
// @Tags Locations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param payload body LocationPayloadDoc true "Location payload"
// @Success 200 {object} LocationResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /locations/{id} [put]
func (h *LocationHandler) updateLocation(w http.ResponseWriter, r *http.Request, id uint) {
	location, err := h.service.GetLocationByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "Lokasi tidak ditemukan"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(location); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}
	if location.Type == "" {
		location.Type = "outlet"
	}

	location.ID = id
	if err := h.service.UpdateLocation(location); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui lokasi"})
		return
	}

	h.writeJSON(w, http.StatusOK, location)
}

// deleteLocation godoc
// @Summary Delete location
// @Description Delete a location record.
// @Tags Locations
// @Security BearerAuth
// @Produce json
// @Param id path int true "Location ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /locations/{id} [delete]
func (h *LocationHandler) deleteLocation(w http.ResponseWriter, r *http.Request, id uint) {
	if err := h.service.DeleteLocation(id); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus lokasi"})
		return
	}

	h.writeJSON(w, http.StatusNoContent, nil)
}

func (h *LocationHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
