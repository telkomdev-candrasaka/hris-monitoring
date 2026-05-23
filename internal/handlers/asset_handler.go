package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/middlewares"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/services"
)

type AssetHandler struct {
	service *services.AssetService
}

type createAssetRequest struct {
	Name         string `json:"name"`
	Category     string `json:"category"`
	SerialNumber string `json:"serial_number"`
	Condition    string `json:"condition,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type borrowAssetRequest struct {
	UserID uint `json:"user_id"`
}

type returnAssetRequest struct {
	Condition string `json:"condition"`
}

func NewAssetHandler(service *services.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

// getAllAssets godoc
// @Summary List assets
// @Description List operational assets and PPE items including status, condition, and borrower linkage.
// @Tags Assets
// @Security BearerAuth
// @Produce json
// @Success 200 {array} AssetResponseDoc
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /assets [get]
func (h *AssetHandler) AssetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllAssets(w, r)
	case http.MethodPost:
		h.createAsset(w, r)
	default:
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

// getAssetByID godoc
// @Summary Get asset by ID
// @Description Retrieve one asset or PPE record by identifier.
// @Tags Assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} AssetResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /assets/{id} [get]
func (h *AssetHandler) AssetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r.URL.Path)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID aset tidak valid"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAssetByID(w, r, id)
	case http.MethodPut:
		h.updateAsset(w, r, id)
	case http.MethodDelete:
		h.deleteAsset(w, r, id)
	default:
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

// BorrowAssetHandler godoc
// @Summary Borrow asset
// @Description Assign an available asset to a user.
// @Tags Assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param payload body BorrowAssetRequestDoc true "Borrow asset payload"
// @Success 200 {object} AssetResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /assets/borrow/{id} [post]
func (h *AssetHandler) BorrowAssetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	id, err := h.parseID(r.URL.Path, "/assets/borrow/")
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID aset tidak valid"})
		return
	}

	var req borrowAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}

	asset, err := h.service.BorrowAsset(id, req.UserID)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusOK, asset)
}

// ReturnAssetHandler godoc
// @Summary Return asset
// @Description Return a borrowed asset and update its latest condition.
// @Tags Assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param payload body ReturnAssetRequestDoc true "Return asset payload"
// @Success 200 {object} AssetResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /assets/return/{id} [post]
func (h *AssetHandler) ReturnAssetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	id, err := h.parseID(r.URL.Path, "/assets/return/")
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID aset tidak valid"})
		return
	}

	var req returnAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}

	asset, err := h.service.ReturnAsset(id, req.Condition)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusOK, asset)
}

// createAsset godoc
// @Summary Create asset
// @Description Create a new PPE or operational asset record.
// @Tags Assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CreateAssetRequestDoc true "Create asset payload"
// @Success 201 {object} AssetResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /assets [post]
func (h *AssetHandler) createAsset(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.ContextUserKey).(*middlewares.JWTClaims)
	if !ok || claims == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token tidak valid"})
		return
	}

	var req createAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}

	asset, err := h.service.CreateAsset(req.Name, req.Category, req.SerialNumber, req.Condition, req.Notes)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusCreated, asset)
}

func (h *AssetHandler) getAllAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := h.service.GetAllAssets()
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data aset"})
		return
	}

	h.writeJSON(w, http.StatusOK, assets)
}

func (h *AssetHandler) getAssetByID(w http.ResponseWriter, r *http.Request, id uint) {
	asset, err := h.service.GetAssetByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "Aset tidak ditemukan"})
		return
	}

	h.writeJSON(w, http.StatusOK, asset)
}

// updateAsset godoc
// @Summary Update asset
// @Description Update editable asset fields such as name, category, condition, and notes.
// @Tags Assets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param payload body UpdateAssetRequestDoc true "Update asset payload"
// @Success 200 {object} AssetResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /assets/{id} [put]
func (h *AssetHandler) updateAsset(w http.ResponseWriter, r *http.Request, id uint) {
	asset, err := h.service.GetAssetByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "Aset tidak ditemukan"})
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}

	if name, ok := req["name"].(string); ok && name != "" {
		asset.Name = name
	}
	if category, ok := req["category"].(string); ok && category != "" {
		asset.Category = category
	}
	if condition, ok := req["condition"].(string); ok && condition != "" {
		asset.Condition = condition
	}
	if notes, ok := req["notes"].(string); ok && notes != "" {
		asset.Notes = notes
	}

	if err := h.service.UpdateAsset(asset); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui aset"})
		return
	}

	h.writeJSON(w, http.StatusOK, asset)
}

// deleteAsset godoc
// @Summary Delete asset
// @Description Delete an asset record.
// @Tags Assets
// @Security BearerAuth
// @Produce json
// @Param id path int true "Asset ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /assets/{id} [delete]
func (h *AssetHandler) deleteAsset(w http.ResponseWriter, r *http.Request, id uint) {
	if err := h.service.DeleteAsset(id); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus aset"})
		return
	}

	h.writeJSON(w, http.StatusNoContent, nil)
}

func (h *AssetHandler) parseID(path string, prefix ...string) (uint, error) {
	var trimmed string
	if len(prefix) > 0 {
		trimmed = strings.TrimPrefix(path, prefix[0])
	} else {
		trimmed = strings.TrimPrefix(path, "/assets/")
	}
	value, err := strconv.Atoi(trimmed)
	return uint(value), err
}

func (h *AssetHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
