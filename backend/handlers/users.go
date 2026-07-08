package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"backend/storage"
)

type UserHandlers struct {
	store *storage.Store
}

func NewUserHandlers(store *storage.Store) *UserHandlers {
	return &UserHandlers{store: store}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if err := h.store.CreateRequest(r.Context(), req.Username, string(bytes)); err != nil {
		http.Error(w, "could not create request", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandlers) GetMe(w http.ResponseWriter, r *http.Request) {
	// TODO: get current user id (from session — later), call h.store.GetUserByID, write JSON response
}
