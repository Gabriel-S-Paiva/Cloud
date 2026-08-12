package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_AllowsUpToLimit(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	handlerCalls := 0
	handler := rl.Limit(func(w http.ResponseWriter, r *http.Request) {
		handlerCalls++
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "203.0.113.5:12345"
		rec := httptest.NewRecorder()
		handler(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d, want %d", i+1, rec.Code, http.StatusOK)
		}
	}
	if handlerCalls != 3 {
		t.Errorf("handlerCalls = %d, want 3", handlerCalls)
	}
}

func TestRateLimiter_BlocksOverLimit(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	handler := rl.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "203.0.113.6:12345"
		rec := httptest.NewRecorder()
		handler(rec, req)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "203.0.113.6:12345"
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("4th request: status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
	if retryAfter := rec.Header().Get("Retry-After"); retryAfter == "" {
		t.Errorf("expected a Retry-After header on a 429 response, got none")
	}
}

func TestRateLimiter_TracksClientsIndependently(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	handler := rl.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req1 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req1.RemoteAddr = "203.0.113.7:1"
	rec1 := httptest.NewRecorder()
	handler(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("client A: status = %d, want %d", rec1.Code, http.StatusOK)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req2.RemoteAddr = "198.51.100.9:1"
	rec2 := httptest.NewRecorder()
	handler(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("client B (different IP): status = %d, want %d", rec2.Code, http.StatusOK)
	}
}

func TestRateLimiter_ResetsAfterWindowElapses(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)
	handler := rl.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "203.0.113.8:1"

	first := httptest.NewRecorder()
	handler(first, req)
	if first.Code != http.StatusOK {
		t.Fatalf("1st request: status = %d, want %d", first.Code, http.StatusOK)
	}

	blocked := httptest.NewRecorder()
	handler(blocked, req)
	if blocked.Code != http.StatusTooManyRequests {
		t.Fatalf("2nd request (immediate): status = %d, want %d", blocked.Code, http.StatusTooManyRequests)
	}

	time.Sleep(60 * time.Millisecond)

	allowed := httptest.NewRecorder()
	handler(allowed, req)
	if allowed.Code != http.StatusOK {
		t.Errorf("3rd request (after window elapsed): status = %d, want %d", allowed.Code, http.StatusOK)
	}
}

func TestRateLimiter_UsesXForwardedForWhenPresent(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	handler := rl.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req1 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req1.RemoteAddr = "127.0.0.1:1"
	req1.Header.Set("X-Forwarded-For", "203.0.113.10")
	rec1 := httptest.NewRecorder()
	handler(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("client A: status = %d, want %d", rec1.Code, http.StatusOK)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req2.RemoteAddr = "127.0.0.1:1"
	req2.Header.Set("X-Forwarded-For", "203.0.113.11")
	rec2 := httptest.NewRecorder()
	handler(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("client B (different XFF, same RemoteAddr): status = %d, want %d", rec2.Code, http.StatusOK)
	}
}
