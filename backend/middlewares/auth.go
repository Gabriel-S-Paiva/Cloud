package middlewares

import (
	"backend/storage"
	"backend/utils"
	"context"
	"errors"
	"net/http"
)

type AuthMiddlewares struct {
	store *storage.Store
}
type contextKey string

const userContextKey contextKey = "user"

func NewAuthMiddlewares(store *storage.Store) *AuthMiddlewares {
	return &AuthMiddlewares{store: store}
}

func (m *AuthMiddlewares) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			utils.WriteJSONError(w, "not authenticated", http.StatusUnauthorized)
			return
		}

		session, err := m.store.GetSession(r.Context(), cookie.Value)
		if errors.Is(err, storage.ErrSessionNotFound) {
			utils.WriteJSONError(w, "session not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, storage.ErrExpiredSession) {
			utils.WriteJSONError(w, "session expired", http.StatusUnauthorized)
		}
		if err != nil {
			utils.WriteJSONError(w, "not authenticated", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, session)
		newReq := r.WithContext(ctx)
		next(w, newReq)
	}
}

func (m *AuthMiddlewares) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(userContextKey).(*storage.Session)
		if !ok {
			utils.WriteJSONError(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		if session.Role != "Admin" {
			utils.WriteJSONError(w, "admin access required", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
