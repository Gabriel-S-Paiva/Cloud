package storage

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
type UserSummary struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
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

func (s *Store) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, username, role, quota, quota_used, root_folder From Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Username, &user.Role, &user.Quota, &user.QuotaUsed, &user.RootFolderId); err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userList, nil
}

func (s *Store) GetSharableUsers(ctx context.Context, userId int) ([]UserSummary, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, username From Users WHERE id != ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []UserSummary
	for rows.Next() {
		var user UserSummary
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userList, nil
}

func (s *Store) SeedAdmin(ctx context.Context, username, password string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var userId int64
	err = tx.QueryRowContext(ctx, "INSERT INTO Users (username, hashed_password, role) VALUES (?, ?, 'Admin') RETURNING id", username, string(hashed)).Scan(&userId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil // admin already exists — not an error, just skip
		}
		return err
	}

	var rootFolderId int64
	err = tx.QueryRowContext(ctx, "INSERT INTO Folders (display_name, owned_by) VALUES ('root', ?) RETURNING id", userId).Scan(&rootFolderId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE Users SET root_folder = ? WHERE id = ?", rootFolderId, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
