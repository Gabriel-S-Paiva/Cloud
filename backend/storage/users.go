package storage

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"-"`
	Role         string `json:"role"`
	Quota        int    `json:"quota"`
	QuotaUsed    int    `json:"quotaUsed"`
	RootFolderId int    `json:"rootFolderId"`
}
type Request struct {
	Id       int
	Username string
	Status   string
}

func (s *Store) GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := s.db.QueryRowContext(ctx, "SELECT id, username, role, quota, quota_used, root_folder FROM Users WHERE Users.id = ?", id).Scan(&user.Id, &user.Username, &user.Role, &user.Quota, &user.QuotaUsed, &user.RootFolderId)
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

func (s *Store) CreateRequest(ctx context.Context, username string, password string) (int, error) {
	var userId int
	err := s.db.QueryRowContext(ctx, "INSERT INTO Requests (username, hashed_password) VALUES (?, ?) RETURNING id", username, password).Scan(userId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, ErrUsernameTaken
		}
		return 0, err
	}
	return userId, nil
}

func (s *Store) AproveRequest(ctx context.Context, id int) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var user User
	err = tx.QueryRowContext(ctx, "SELECT username, hashed_password FROM Requests WHERE Requests.id = ?", id).Scan(&user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return 0, ErrUserNotFound
	}
	if err != nil {
		return 0, err
	}

	err = tx.QueryRowContext(ctx, "INSERT INTO Users (username, hashed_password) VALUES (?, ?) RETURNING id", user.Username, user.Password).Scan(&user.Id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, ErrUsernameTaken
		}
		return 0, err
	}

	var rootFolderId int64
	err = tx.QueryRowContext(ctx, "INSERT INTO Folder (display_name, owned_by) VALUES ('root',?) RETURNING id", user.Id).Scan(rootFolderId)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, "UPDATE Users SET root_folder = ? WHERE id = ?", rootFolderId, user.Id)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM Requests WHERE Requests.id = ?", id)
	if err != nil {
		return 0, err
	}
	return user.Id, tx.Commit()
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
