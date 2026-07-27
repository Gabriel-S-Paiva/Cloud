package storage

import (
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const storageDir = "storage"

type File struct {
	Id            int       `json:"id"`
	DisplayName   string    `json:"displayName"`
	OwnedBy       int       `json:"ownedBy"`
	Size          int       `json:"size"`
	BytesReceived int       `json:"bytesReceived"`
	Status        string    `json:"status"`
	ContentType   string    `json:"contentType"`
	UploadedAt    int       `json:"uploadedAt"`
	LastModified  int       `json:"lastModified"`
	ParentFolder  NullInt64 `json:"parentFolder"`
}

func (s *Store) CreateFileUploadIntent(ctx context.Context, displayName string, ownedBy int, parentFolder sql.NullInt64, size int, contentType string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var quota, quotaUsed int
	err = tx.QueryRowContext(ctx, "SELECT quota, quota_used FROM Users WHERE id = ?", ownedBy).
		Scan(&quota, &quotaUsed)
	if err != nil {
		return 0, err
	}

	if quotaUsed+size > quota {
		return 0, ErrQuotaExceeded
	}

	result, err := tx.ExecContext(ctx, "INSERT INTO Files (display_name, owned_by, size, parent_folder, content_type, status, bytes_received) VALUES (?,?,?,?,?,'Uploading',0)", displayName, ownedBy, size, parentFolder, contentType)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	_, err = tx.ExecContext(ctx, "UPDATE Users SET quota_used = quota_used + ? WHERE id = ?", size, ownedBy)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Store) AppendFileChunk(ctx context.Context, fileId int, chunk io.Reader) (int64, error) {
	file, err := s.GetFileById(ctx, fileId)
	if err != nil {
		return 0, err
	}
	if file.Status != "Uploading" {
		return 0, ErrFileNotUploading
	}

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return 0, err
	}
	path := filepath.Join(storageDir, strconv.Itoa(fileId))

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	written, err := io.Copy(f, chunk)
	if err != nil {
		return 0, err
	}

	newTotal := file.BytesReceived + int(written)
	status := "Uploading"
	if newTotal >= file.Size {
		status = "Complete"
	}

	_, err = s.db.ExecContext(ctx, "UPDATE Files SET bytes_received = ?, status = ? WHERE id = ?", newTotal, status, fileId)
	if err != nil {
		return 0, err
	}

	return written, nil
}

func (s *Store) GetFileById(ctx context.Context, id int) (*File, error) {
	var file File
	err := s.db.QueryRowContext(ctx, "SELECT id, display_name, owned_by, size, bytes_received, status, content_type, uploaded_at, last_modified, parent_folder FROM Files WHERE id = ?", id).
		Scan(&file.Id, &file.DisplayName, &file.OwnedBy, &file.Size, &file.BytesReceived, &file.Status, &file.ContentType, &file.UploadedAt, &file.LastModified, &file.ParentFolder)
	if err == sql.ErrNoRows {
		return nil, ErrFileNotFound
	}
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *Store) GetFileContent(ctx context.Context, id int) (*File, io.ReadCloser, error) {
	file, err := s.GetFileById(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if file.Status != "Complete" {
		return nil, nil, ErrFileNotComplete
	}

	path := filepath.Join(storageDir, strconv.Itoa(id))
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	return file, f, nil
}

func (s *Store) UpdateFile(ctx context.Context, id int, newName *string, newParent *sql.NullInt64) error {
	if newParent != nil {
		res, err := s.db.ExecContext(ctx, "UPDATE Files SET parent_folder = ? WHERE id = ?", newParent, id)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows != 1 {
			return ErrUpdatingFile
		}
	}
	if newName != nil {
		res, err := s.db.ExecContext(ctx, "UPDATE Files SET display_name = ? WHERE id = ?", newName, id)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows != 1 {
			return ErrUpdatingFile
		}
	}
	_, err := s.db.ExecContext(ctx, "UPDATE Files SET last_modified = ? WHERE id = ?", time.Now().Unix(), id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteFile(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	var fileSize int64

	err = tx.QueryRowContext(ctx, "DELETE FROM Files WHERE id = ? RETURNING owned_by, size", id).
		Scan(&userID, &fileSize)

	if err == sql.ErrNoRows {
		return ErrFileNotFound
	}
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE Users SET quota_used = quota_used - ? WHERE id = ?", fileSize, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) FileOwnership(ctx context.Context, fileId int, userId int) (*File, error) {
	file, err := s.GetFileById(ctx, fileId)
	if err != nil {
		return nil, err
	}
	if file.OwnedBy != userId {
		return nil, ErrForbidden
	}
	return file, nil
}

func (s *Store) FileAccess(ctx context.Context, fileId, userId int) (*File, string, error) {
	var file File
	var permission sql.NullString
	err := s.db.QueryRowContext(ctx,
		`SELECT f.id, f.display_name, f.owned_by, f.size, f.bytes_received, f.status, f.content_type, f.uploaded_at, f.last_modified, f.parent_folder, s.permissions
		 FROM Files f
		 LEFT JOIN Shares s ON s.file = f.id AND s.shared_with = ?
		 WHERE f.id = ? AND (f.owned_by = ? OR s.shared_with = ?)`,
		userId, fileId, userId, userId).
		Scan(&file.Id, &file.DisplayName, &file.OwnedBy, &file.Size, &file.BytesReceived, &file.Status, &file.ContentType, &file.UploadedAt, &file.LastModified, &file.ParentFolder, &permission)
	if err == sql.ErrNoRows {
		return nil, "", ErrFileNotFound
	}
	if err != nil {
		return nil, "", err
	}

	if file.OwnedBy == userId {
		return &file, "Owner", nil
	}
	return &file, permission.String, nil
}
