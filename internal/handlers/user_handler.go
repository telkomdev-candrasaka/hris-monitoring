package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

type createUserRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	LocationID uint   `json:"location_id"`
	ShiftID    *uint  `json:"shift_id,omitempty"`
	Password   string `json:"password"`
	BaseSalary float64 `json:"base_salary"`
}

type updateUserRequest struct {
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	Role       string `json:"role,omitempty"`
	LocationID uint   `json:"location_id,omitempty"`
	ShiftID    *uint  `json:"shift_id,omitempty"`
	Password   string `json:"password,omitempty"`
	BaseSalary float64 `json:"base_salary,omitempty"`
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes() {
	http.HandleFunc("/users", h.usersHandler)
	http.HandleFunc("/users/", h.userByIDHandler)
}

func (h *UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	h.usersHandler(w, r)
}

func (h *UserHandler) UserByIDHandler(w http.ResponseWriter, r *http.Request) {
	h.userByIDHandler(w, r)
}

// getAllUsers godoc
// @Summary List users
// @Description List employees and managers with their assigned location, shift, and base salary.
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} UserResponseDoc
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (h *UserHandler) userByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r.URL.Path)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID user tidak valid"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUserByID(w, r, id)
	case http.MethodPut:
		h.updateUser(w, r, id)
	case http.MethodDelete:
		h.deleteUser(w, r, id)
	default:
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (h *UserHandler) parseID(path string) (uint, error) {
	trimmed := strings.TrimPrefix(path, "/users/")
	value, err := strconv.Atoi(trimmed)
	return uint(value), err
}

// createUser godoc
// @Summary Create user
// @Description Create a new user account and assign role, location, optional shift, and base salary.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CreateUserRequestDoc true "User payload"
// @Success 201 {object} UserResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}

	user := models.User{
		Name:       req.Name,
		Email:      req.Email,
		Role:       req.Role,
		LocationID: req.LocationID,
		ShiftID:    req.ShiftID,
		Password:   req.Password,
		BaseSalary: req.BaseSalary,
	}

	if err := h.service.CreateUser(&user); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal membuat user"})
		return
	}

	h.writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil daftar user"})
		return
	}

	h.writeJSON(w, http.StatusOK, users)
}

// getUserByID godoc
// @Summary Get user by ID
// @Description Retrieve one user by identifier.
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request, id uint) {
	user, err := h.service.GetUserByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "User tidak ditemukan"})
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// updateUser godoc
// @Summary Update user
// @Description Update mutable user fields such as role, location, shift, password, and salary.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param payload body UpdateUserRequestDoc true "Update user payload"
// @Success 200 {object} UserResponseDoc
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id uint) {
	user, err := h.service.GetUserByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "User tidak ditemukan"})
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Payload tidak valid"})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.LocationID != 0 {
		user.LocationID = req.LocationID
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.ShiftID != nil {
		user.ShiftID = req.ShiftID
	}
	if req.BaseSalary != 0 {
		user.BaseSalary = req.BaseSalary
	}

	if err := h.service.UpdateUser(user); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui user"})
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// deleteUser godoc
// @Summary Delete user
// @Description Delete a user. Deletion is blocked if the user still holds borrowed assets.
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, id uint) {
	if err := h.service.DeleteUser(id); err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus user"})
		return
	}

	h.writeJSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
