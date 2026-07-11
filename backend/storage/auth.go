package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"time"
)

type Session struct {
	Token     string
	UserId    int
	ExpiresAt int64
	Role      string
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := s.db.QueryRowContext(ctx, "SELECT id, username, hashed_password, role, quota, quota_used FROM Users WHERE username = ?", username).Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.Quota, &user.QuotaUsed)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrUserNotFound
	case err != nil:
		return nil, err
	default:
		return &user, nil
	}
}

func (s *Store) CreateSession(ctx context.Context, userID int, expiresAt time.Time) (string, error) {
	token := rand.Text()

	_, err := s.db.ExecContext(ctx, "INSERT INTO Sessions (token, user_id, expires_at) VALUES (?, ?, ?)", token, userID, expiresAt.Unix())
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Store) GetSession(ctx context.Context, token string) (*Session, error) {
	var session Session
	err := s.db.QueryRowContext(ctx, "SELECT token, user_id, expires_at, role FROM Sessions s JOIN Users u ON s.user_id = u.id WHERE s.token = ?", token).Scan(&session.Token, &session.UserId, &session.ExpiresAt, &session.Role)
	if err == sql.ErrNoRows {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	if time.Time.Unix(time.Now()) > session.ExpiresAt {
		return nil, ErrExpiredSession
	}

	return &session, nil
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM Sessions WHERE token=?", token)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSessionNotFound
	}
	if rows != 1 {
		return ErrSessionRowMismatch
	}
	return nil
}
