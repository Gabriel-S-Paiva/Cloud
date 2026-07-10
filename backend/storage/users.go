package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
)

var ErrUsernameTaken = errors.New("Username already taken.")
var ErrUserNotFound = errors.New("Username not found")

type User struct {
	Id        int
	Username  string
	Role      string
	Quota     int
	QuotaUsed int
}

func (s *Store) GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := s.db.QueryRowContext(ctx, "SELECT id, username, role, quota, quota_used FROM Users WHERE Users.id = ?", id).Scan(&user.Id, &user.Username, &user.Role, &user.Quota, &user.QuotaUsed)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("User %v not found", id)
		return nil, ErrUserNotFound
	case err != nil:
		return nil, err
	default:
		return &user, nil
	}
}

func (s *Store) CreateRequest(ctx context.Context, username string, password string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO Requests (username, hashed_password) VALUES (?, ?)", username, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUsernameTaken
		}
		return err
	}
	return nil
}
