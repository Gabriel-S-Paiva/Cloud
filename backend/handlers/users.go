package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"backend/storage"
)

type UserHandlers struct {
	store *storage.Store
}
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserHandlers(store *storage.Store) *UserHandlers {
	return &UserHandlers{store: store}
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		writeJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandlers) AproveRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, "Error parsing id", http.StatusBadRequest)
		return
	}

	err = h.store.AproveRequest(r.Context(), id)
	if errors.Is(err, storage.ErrUsernameTaken) {
		writeJSONError(w, "Username already taken", http.StatusConflict)
		return
	}
	if err != nil {
		writeJSONError(w, "Error aproving Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandlers) RejectRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, fmt.Sprintf("Error parsing %v", id), http.StatusBadRequest)
		return
	}

	err = h.store.RejectRequest(r.Context(), id)
	if errors.Is(err, storage.ErrUpdatingRequest) {
		writeJSONError(w, "Error Updating Request", http.StatusInternalServerError)
		return
	}
	if err != nil {
		writeJSONError(w, "Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandlers) ListPendingRequests(w http.ResponseWriter, r *http.Request) {
	requestsList, err := h.store.ListPendingRequests(r.Context())
	if err != nil {
		writeJSONError(w, "Error Searching Table", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(requestsList); err != nil {
		writeJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}
