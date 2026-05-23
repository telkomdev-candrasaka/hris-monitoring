package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes() {
	http.HandleFunc("/login", h.login)
}

// login godoc
// @Summary Login user
// @Description Authenticate user with email and password, then return a JWT bearer token for protected endpoints.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body LoginRequestDoc true "Login payload"
// @Success 200 {object} LoginResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 405 {object} ErrorResponse
// @Example request {"email":"admin@example.com","password":"secret123"}
// @Router /login [post]
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Email atau password salah"})
		return
	}

	h.writeJSON(w, http.StatusOK, loginResponse{Token: token})
}

func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
