package storage

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

type User struct {
	Id        int
	Username  string
	Password  string
	Role      string
	Quota     int
	QuotaUsed int
}
type Request struct {
	Id       int
	Username string
	Status   string
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

func (s *Store) AproveRequest(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var user User
	err = tx.QueryRowContext(ctx, "SELECT username, hashed_password FROM Requests WHERE Requests.id = ?", id).Scan(&user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO Users (username, hashed_password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUsernameTaken
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM Requests WHERE Requests.id = ?", id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) RejectRequest(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, "UPDATE Requests SET status = 'Rejected' WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrUpdatingRequest
	}
	return nil
}

func (s *Store) ListPendingRequests(ctx context.Context) ([]Request, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, username, status FROM Requests WHERE status = 'Pending'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requestList []Request
	for rows.Next() {
		var request Request
		if err := rows.Scan(&request.Id, &request.Username, &request.Status); err != nil {
			return nil, err
		}
		requestList = append(requestList, request)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requestList, nil
}
