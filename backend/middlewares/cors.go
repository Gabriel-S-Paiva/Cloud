package middlewares

import (
	"net/http"
	"os"
	"strings"
)

var defaultOrigins = []string{
	"http://localhost:5173",
	"http://localhost:3000",
}

func CORS(next http.Handler) http.Handler {
	allowedMap := make(map[string]bool)

	rawEnv := os.Getenv("ALLOWED_ORIGINS")
	if rawEnv != "" {
		for _, origin := range strings.Split(rawEnv, ",") {
			allowedMap[strings.TrimSpace(origin)] = true
		}
	} else {
		for _, origin := range defaultOrigins {
			allowedMap[origin] = true
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Only set origin & credentials if the requesting origin is explicitly allowed
		if allowedMap[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
