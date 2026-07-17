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
	result, err := s.db.ExecContext(ctx, "INSERT INTO Files (display_name, owned_by, size, parent_folder) VALUES (?,?,?,?)", displayName, ownedBy, size, parentFolder)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
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
	result, err := s.db.ExecContext(ctx, "DELETE FROM Files WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrFileNotFound
	}
	if rows != 1 {
		return ErrFileMismatch
	}
	return nil
}
