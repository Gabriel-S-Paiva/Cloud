package storage

import (
	"context"
	"database/sql"
	"time"
)

type File struct {
	Id           int
	DisplayName  string
	OwnedBy      int
	Size         int
	UploadedAt   int
	LastModified int
	ParentFolder sql.NullInt64
}

func (s *Store) CreateFile(ctx context.Context, displayName string, ownedBy int, parentFolder sql.NullInt64, size int) (int64, error) {
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

	result, err := tx.ExecContext(ctx, "INSERT INTO Files (display_name, owned_by, size, parent_folder) VALUES (?,?,?,?)", displayName, ownedBy, size, parentFolder)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	_, err = tx.ExecContext(ctx, "UPDATE Users SET quota_used VALUES quota_used + ? WHERE id = ?", size, ownedBy)
	tx.Commit()

	return id, nil
}

func (s *Store) GetFileById(ctx context.Context, id int) (*File, error) {
	var file File
	err := s.db.QueryRowContext(ctx, "SELECT id, display_name, owned_by, size, uploaded_at, last_modified, parent_folder FROM Files WHERE id = ?", id).Scan(&file.Id, &file.DisplayName, &file.OwnedBy, &file.Size, &file.UploadedAt, &file.LastModified, &file.ParentFolder)
	if err == sql.ErrNoRows {
		return nil, ErrFileNotFound
	}
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *Store) GetFile(ctx context.Context, id int) {

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
