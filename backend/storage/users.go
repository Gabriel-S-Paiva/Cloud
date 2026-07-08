package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

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
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &user, nil
	}
}

func (s *Store) CreateRequest(ctx context.Context, username string, password string) error {
	result, err := s.db.ExecContext(ctx, "INSERT INTO Requests (username, hashed_password) VALUES (?, ?)", username, password)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		err = fmt.Errorf("expected to affect 1 row, affected %v", rows)
		return err
	}
	return nil
}
