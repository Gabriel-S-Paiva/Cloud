package handlers

import (
	"backend/storage"
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
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.store.GetUserByUsername(r.Context(), loginForm.Username)
	if errors.Is(err, storage.ErrUserNotFound) {
		writeJSONError(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}
	if err != nil {
		writeJSONError(w, "could not create request", http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password))
	if err != nil {
		writeJSONError(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}

	sessionDuration := 24 * time.Hour
	expiresAt := time.Now().Add(sessionDuration)

	token, err := h.store.CreateSession(r.Context(), user.Id, expiresAt)
	if err != nil {
		writeJSONError(w, "Could not create session", http.StatusInternalServerError)
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
		writeJSONError(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	if err := h.store.DeleteSession(r.Context(), cookie.Value); err != nil {
		writeJSONError(w, "Error Ending Session", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
