package handlers

import (
	"backend/storage"
	"backend/utils"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	store *storage.Store
}
type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthHandlers(store *storage.Store) *AuthHandlers {
	return &AuthHandlers{store: store}
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var loginForm LoginForm
	if err := json.NewDecoder(r.Body).Decode(&loginForm); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.store.GetUserByUsername(r.Context(), loginForm.Username)
	if errors.Is(err, storage.ErrUserNotFound) {
		utils.WriteJSONError(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "could not create request", http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password))
	if err != nil {
		utils.WriteJSONError(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}

	sessionDuration := 24 * time.Hour
	expiresAt := time.Now().Add(sessionDuration)

	token, err := h.store.CreateSession(r.Context(), user.Id, expiresAt)
	if err != nil {
		utils.WriteJSONError(w, "Could not create session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false, // DONT FORGET CHANGE TO TRUE AFTER LEAVING LOCALHOST
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJSONError(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	if err := h.store.DeleteSession(r.Context(), cookie.Value); err != nil {
		utils.WriteJSONError(w, "Error Ending Session", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
