package storage

import (
	"context"
	"testing"
	"time"
)

func TestGetUserByUsername_Success(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	seedUser(t, s, "authlookup", 1000)

	user, err := s.GetUserByUsername(ctx, "authlookup")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if user.Username != "authlookup" {
		t.Errorf("Username = %q, want %q", user.Username, "authlookup")
	}
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetUserByUsername(ctx, "ghost")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestCreateSession_ThenGetSession(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "sessionuser", 1000)

	expiresAt := time.Now().Add(time.Hour)
	token, err := s.CreateSession(ctx, userID, expiresAt)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if token == "" {
		t.Fatal("CreateSession returned an empty token")
	}

	session, err := s.GetSession(ctx, token)
	if err != nil {
		t.Fatalf("GetSession: %v", err)
	}
	if session.UserId != userID {
		t.Errorf("UserId = %d, want %d", session.UserId, userID)
	}
	if session.Role != "User" {
		t.Errorf("Role = %q, want %q", session.Role, "User")
	}
}

func TestGetSession_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetSession(ctx, "nonexistent-token")
	if err != ErrSessionNotFound {
		t.Errorf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestGetSession_Expired(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "expireduser", 1000)

	expiresAt := time.Now().Add(-time.Hour) // already in the past
	token, err := s.CreateSession(ctx, userID, expiresAt)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	_, err = s.GetSession(ctx, token)
	if err != ErrExpiredSession {
		t.Errorf("expected ErrExpiredSession, got %v", err)
	}
}

func TestDeleteSession_Success(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "logoutuser", 1000)

	token, err := s.CreateSession(ctx, userID, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	if err := s.DeleteSession(ctx, token); err != nil {
		t.Fatalf("DeleteSession: %v", err)
	}
	if _, err := s.GetSession(ctx, token); err != ErrSessionNotFound {
		t.Errorf("expected ErrSessionNotFound after delete, got %v", err)
	}
}

func TestDeleteSession_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.DeleteSession(ctx, "never-existed")
	if err != ErrSessionNotFound {
		t.Errorf("expected ErrSessionNotFound, got %v", err)
	}
}
