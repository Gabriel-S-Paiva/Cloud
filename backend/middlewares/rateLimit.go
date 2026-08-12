package middlewares

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"backend/utils"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := clientKey(r)
		if !rl.allow(key) {
			w.Header().Set("Retry-After", strconv.Itoa(int(rl.window.Seconds())))
			utils.WriteJSONError(w, "too many requests, please try again later", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func (rl *RateLimiter) allow(key string) bool {
	now := time.Now()
	cutoff := now.Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	existing := rl.requests[key]
	kept := existing[:0]
	for _, t := range existing {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}

	if len(kept) >= rl.limit {
		rl.requests[key] = kept
		return false
	}

	rl.requests[key] = append(kept, now)
	return true
}

func clientKey(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if idx := strings.IndexByte(fwd, ','); idx != -1 {
			return strings.TrimSpace(fwd[:idx])
		}
		return strings.TrimSpace(fwd)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
