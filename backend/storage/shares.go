package storage

import (
	"context"
	"database/sql"
)

type Share struct {
	Id         int       `json:"id"`
	File       NullInt64 `json:"file"`
	Folder     NullInt64 `json:"folder"`
	SharedWith int       `json:"sharedWith"`
	Permission string    `json:"permission"`
}

func (s *Store) CreateShare(ctx context.Context, file sql.NullInt64, folder sql.NullInt64, sharedWith int, permission string) (int, error) {
	var shareId int
	if err := s.db.QueryRowContext(ctx, "INSERT INTO Shares (file, folder, shared_with, permissions) VALUES (?,?,?,?) RETURNING id", file, folder, sharedWith, permission).Scan(shareId); err != nil {
		return 0, err
	}
	return shareId, nil
}

func (s *Store) GetShareById(ctx context.Context, id int) (*Share, error) {
	var share Share
	err := s.db.QueryRowContext(ctx, "SELECT id, file, folder, shared_with, permission FROM Shares WHERE id = ?", id).Scan(&share.Id, &share.File, &share.Folder, &share.SharedWith, &share.Permission)
	if err == sql.ErrNoRows {
		return nil, ErrShareNotFound
	}
	if err != nil {
		return nil, err
	}
	return &share, nil
}

// Shared With Me
func (s *Store) GetIncomingShares(ctx context.Context, userId int) (*FolderContents, error) {
	var folderContents FolderContents
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	fileRows, err := tx.QueryContext(ctx, "SELECT id, display_name, owned_by, size, uploaded_at, last_modified, parent_folder FROM Files AS f JOIN Shares AS s ON f.id = s.file WHERE s.shared_with = ?", userId)
	if err != nil {
		return nil, err
	}
	defer fileRows.Close()
	for fileRows.Next() {
		var file File
		if err := fileRows.Scan(&file.Id, &file.DisplayName, &file.OwnedBy, &file.Size, &file.UploadedAt, &file.LastModified, &file.ParentFolder); err != nil {
			return nil, err
		}
		folderContents.Files = append(folderContents.Files, file)
	}
	if err := fileRows.Err(); err != nil {
		return nil, err
	}

	folderRows, err := tx.QueryContext(ctx, "SELECT id, display_name, owned_by, parent_folder FROM Folders AS f JOIN Shares AS s ON f.id = s.folder WHERE s.shared_with = ?", userId)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()
	for folderRows.Next() {
		var folder Folder
		if err := folderRows.Scan(&folder.Id, &folder.DisplayName, &folder.OwnedBy, &folder.ParentFolder); err != nil {
			return nil, err
		}
		folderContents.Folders = append(folderContents.Folders, folder)
	}
	if err := folderRows.Err(); err != nil {
		return nil, err
	}
	return &folderContents, nil
}

// Shared Files
func (s *Store) GetOutgoingShares(ctx context.Context, userId int) (*FolderContents, error) {
	var folderContents FolderContents
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	fileRows, err := tx.QueryContext(ctx, "SELECT id, display_name, owned_by, size, uploaded_at, last_modified, parent_folder FROM Files AS f JOIN Shares AS s ON f.id = s.file WHERE f.owned_by = ?", userId)
	if err != nil {
		return nil, err
	}
	defer fileRows.Close()
	for fileRows.Next() {
		var file File
		if err := fileRows.Scan(&file.Id, &file.DisplayName, &file.OwnedBy, &file.Size, &file.UploadedAt, &file.LastModified, &file.ParentFolder); err != nil {
			return nil, err
		}
		folderContents.Files = append(folderContents.Files, file)
	}
	if err := fileRows.Err(); err != nil {
		return nil, err
	}

	folderRows, err := tx.QueryContext(ctx, "SELECT id, display_name, owned_by, parent_folder FROM Folders AS f JOIN Shares AS s ON f.id = s.folder WHERE f.owned_by = ?", userId)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()
	for folderRows.Next() {
		var folder Folder
		if err := folderRows.Scan(&folder.Id, &folder.DisplayName, &folder.OwnedBy, &folder.ParentFolder); err != nil {
			return nil, err
		}
		folderContents.Folders = append(folderContents.Folders, folder)
	}
	if err := folderRows.Err(); err != nil {
		return nil, err
	}
	return &folderContents, nil
}

func (s *Store) UpdatePermission(ctx context.Context, shareId int, newPermission string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE Shares SET permissions = ? WHERE id = ?", newPermission, shareId)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrUpdatingShare
	}
	return nil
}

func (s *Store) DeleteShare(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM Shares WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrShareNotFound
	}
	return nil
}
