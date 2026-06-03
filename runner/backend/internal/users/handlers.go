package users

import (
	"net/http"

	"github.com/quoromx/irondoj/backend/internal/auth"
	"github.com/quoromx/irondoj/backend/internal/database"
	"github.com/quoromx/irondoj/backend/internal/platform"
)

type Handlers struct {
	store *database.Store
	auth  *auth.Service
}

func New(store *database.Store, authService *auth.Service) *Handlers {
	return &Handlers{store: store, auth: authService}
}

type authRequest struct {
	Handle      string `json:"handle"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) error {
	var req authRequest
	if err := platform.Decode(r, &req); err != nil {
		return platform.Error{Status: http.StatusBadRequest, Message: "invalid registration payload"}
	}
	user, token, err := h.auth.Register(req.Handle, req.DisplayName, req.Email, req.Password)
	if err != nil {
		return platform.Error{Status: http.StatusBadRequest, Message: err.Error()}
	}
	return platform.JSON(w, http.StatusCreated, map[string]any{"user": user, "token": token})
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) error {
	var req authRequest
	if err := platform.Decode(r, &req); err != nil {
		return platform.Error{Status: http.StatusBadRequest, Message: "invalid login payload"}
	}
	user, token, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		return platform.Error{Status: http.StatusUnauthorized, Message: err.Error()}
	}
	return platform.JSON(w, http.StatusOK, map[string]any{"user": user, "token": token})
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) error {
	user, ok := h.auth.UserFromToken(platform.BearerToken(r))
	if !ok {
		return platform.Error{Status: http.StatusUnauthorized, Message: "authentication required"}
	}
	return platform.JSON(w, http.StatusOK, map[string]any{"user": user})
}

func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) error {
	user, ok := h.auth.UserFromToken(platform.BearerToken(r))
	if !ok {
		return platform.Error{Status: http.StatusUnauthorized, Message: "authentication required"}
	}
	dashboard, ok := h.store.Dashboard(user.ID)
	if !ok {
		return platform.Error{Status: http.StatusNotFound, Message: "dashboard not found"}
	}
	return platform.JSON(w, http.StatusOK, dashboard)
}

func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) error {
	user, ok := h.auth.UserFromToken(platform.BearerToken(r))
	if !ok {
		return platform.Error{Status: http.StatusUnauthorized, Message: "authentication required"}
	}
	return platform.JSON(w, http.StatusOK, map[string]any{"user": user, "token": h.auth.Token(user.ID)})
}

func (h *Handlers) Logout(w http.ResponseWriter, _ *http.Request) error {
	return platform.JSON(w, http.StatusOK, map[string]any{"message": "logged out"})
}
