package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"backend/middlewares"
	"backend/storage"
	"backend/utils"
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
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJSONError(w, "internal error", http.StatusInternalServerError)
		return
	}
	id, err := h.store.CreateRequest(r.Context(), req.Username, string(bytes))
	if errors.Is(err, storage.ErrUsernameTaken) {
		utils.WriteJSONError(w, "username already taken", http.StatusConflict)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "could not create request", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *UserHandlers) GetMe(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := h.store.GetUserByID(r.Context(), session.UserId)
	if errors.Is(err, storage.ErrUserNotFound) {
		utils.WriteJSONError(w, fmt.Sprintf("User %v not found", session.UserId), http.StatusNotFound)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Error Searching User", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		utils.WriteJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandlers) AproveRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "Error parsing id", http.StatusBadRequest)
		return
	}

	userId, err := h.store.AproveRequest(r.Context(), id)
	if errors.Is(err, storage.ErrUsernameTaken) {
		utils.WriteJSONError(w, "Username already taken", http.StatusConflict)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Error aproving Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": userId})
}

func (h *UserHandlers) RejectRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, fmt.Sprintf("Error parsing %v", id), http.StatusBadRequest)
		return
	}

	err = h.store.RejectRequest(r.Context(), id)
	if errors.Is(err, storage.ErrUpdatingRequest) {
		utils.WriteJSONError(w, "Error Updating Request", http.StatusInternalServerError)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandlers) ListPendingRequests(w http.ResponseWriter, r *http.Request) {
	requestsList, err := h.store.ListPendingRequests(r.Context())
	if err != nil {
		utils.WriteJSONError(w, "Error Searching Table", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(requestsList); err != nil {
		utils.WriteJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}
