package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, "internal error", http.StatusInternalServerError)
		return
	}
	err = h.store.CreateRequest(r.Context(), req.Username, string(bytes))
	if errors.Is(err, storage.ErrUsernameTaken) {
		writeJSONError(w, "username already taken", http.StatusConflict)
		return
	}
	if err != nil {
		writeJSONError(w, "could not create request", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandlers) GetMe(w http.ResponseWriter, r *http.Request) {
	id := 1 //TODO: add JWT for the me id

	user, err := h.store.GetUserByID(r.Context(), id)
	if errors.Is(err, storage.ErrUserNotFound) {
		writeJSONError(w, fmt.Sprintf("User %v not found", id), http.StatusNotFound)
		return
	}
	if err != nil {
		writeJSONError(w, "Error Searching User", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		writeJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}
