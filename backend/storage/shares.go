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

type SharedFolder struct {
	Folder
	ShareId         int    `json:"shareId"`
	SharedWith      string `json:"sharedWith"`
	OwnedByUsername string `json:"ownedByUsername"`
	Permissions     string `json:"permissions"`
}

type SharedFile struct {
	File
	ShareId         int    `json:"shareId"`
	SharedWith      string `json:"sharedWith"`
	OwnedByUsername string `json:"ownedByUsername"`
	Permissions     string `json:"permissions"`
}

type SharedContents struct {
	Folders []SharedFolder `json:"folders"`
	Files   []SharedFile   `json:"files"`
}

func (s *Store) CreateShare(ctx context.Context, file sql.NullInt64, folder sql.NullInt64, sharedWith int, permission string) (int, error) {
	var shareId int
	if err := s.db.QueryRowContext(ctx, "INSERT INTO Shares (file, folder, shared_with, permissions) VALUES (?,?,?,?) RETURNING id", file, folder, sharedWith, permission).Scan(&shareId); err != nil {
		return 0, err
	}
	return shareId, nil
}

func (s *Store) GetShareById(ctx context.Context, id int) (*Share, error) {
	var share Share
	err := s.db.QueryRowContext(ctx, "SELECT id, file, folder, shared_with, permissions FROM Shares WHERE id = ?", id).Scan(&share.Id, &share.File, &share.Folder, &share.SharedWith, &share.Permission)
	if err == sql.ErrNoRows {
		return nil, ErrShareNotFound
	}
	if err != nil {
		return nil, err
	}
	return &share, nil
}

// Shared With Me
func (s *Store) GetIncomingShares(ctx context.Context, userId int) (*SharedContents, error) {
	var contents SharedContents
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	fileRows, err := tx.QueryContext(ctx, `SELECT f.id, f.display_name, f.owned_by, f.size, f.bytes_received, f.status, f.content_type, f.uploaded_at, f.last_modified, f.parent_folder, s.id, s.permissions, r.username, owner.username
											FROM Files AS f
											JOIN Shares AS s ON f.id = s.file
											JOIN Users AS r ON s.shared_with = r.id
											JOIN Users AS owner ON f.owned_by = owner.id
											WHERE s.shared_with = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer fileRows.Close()
	for fileRows.Next() {
		var sf SharedFile
		if err := fileRows.Scan(&sf.Id, &sf.DisplayName, &sf.OwnedBy, &sf.Size, &sf.BytesReceived, &sf.Status, &sf.ContentType, &sf.UploadedAt, &sf.LastModified, &sf.ParentFolder, &sf.ShareId, &sf.Permissions, &sf.SharedWith, &sf.OwnedByUsername); err != nil {
			return nil, err
		}
		contents.Files = append(contents.Files, sf)
	}
	if err := fileRows.Err(); err != nil {
		return nil, err
	}

	folderRows, err := tx.QueryContext(ctx, `SELECT f.id, f.display_name, f.owned_by, f.parent_folder, s.id, s.permissions, r.username, owner.username
												FROM Folders AS f
												JOIN Shares AS s ON f.id = s.folder
												JOIN Users AS r ON s.shared_with = r.id
												JOIN Users AS owner ON f.owned_by = owner.id
												WHERE s.shared_with = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()
	for folderRows.Next() {
		var sf SharedFolder
		if err := folderRows.Scan(&sf.Id, &sf.DisplayName, &sf.OwnedBy, &sf.ParentFolder, &sf.ShareId, &sf.Permissions, &sf.SharedWith, &sf.OwnedByUsername); err != nil {
			return nil, err
		}
		contents.Folders = append(contents.Folders, sf)
	}
	if err := folderRows.Err(); err != nil {
		return nil, err
	}
	return &contents, nil
}

// Shared Files
func (s *Store) GetOutgoingShares(ctx context.Context, userId int) (*SharedContents, error) {
	var contents SharedContents
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	fileRows, err := tx.QueryContext(ctx, `SELECT f.id, f.display_name, f.owned_by, f.size, f.bytes_received, f.status, f.content_type, f.uploaded_at, f.last_modified, f.parent_folder, s.id, s.permissions, r.username, owner.username
											FROM Files AS f
											JOIN Shares AS s ON f.id = s.file
											JOIN Users AS r ON s.shared_with = r.id
											JOIN Users AS owner ON f.owned_by = owner.id
											WHERE f.owned_by = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer fileRows.Close()
	for fileRows.Next() {
		var sf SharedFile
		if err := fileRows.Scan(&sf.Id, &sf.DisplayName, &sf.OwnedBy, &sf.Size, &sf.BytesReceived, &sf.Status, &sf.ContentType, &sf.UploadedAt, &sf.LastModified, &sf.ParentFolder, &sf.ShareId, &sf.Permissions, &sf.SharedWith, &sf.OwnedByUsername); err != nil {
			return nil, err
		}
		contents.Files = append(contents.Files, sf)
	}
	if err := fileRows.Err(); err != nil {
		return nil, err
	}

	folderRows, err := tx.QueryContext(ctx, `SELECT f.id, f.display_name, f.owned_by, f.parent_folder, s.id, s.permissions, r.username, owner.username
												FROM Folders AS f
												JOIN Shares AS s ON f.id = s.folder
												JOIN Users AS r ON r.id = s.shared_with
												JOIN Users AS owner ON f.owned_by = owner.id
												WHERE f.owned_by = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()
	for folderRows.Next() {
		var sf SharedFolder
		if err := folderRows.Scan(&sf.Id, &sf.DisplayName, &sf.OwnedBy, &sf.ParentFolder, &sf.ShareId, &sf.Permissions, &sf.SharedWith, &sf.OwnedByUsername); err != nil {
			return nil, err
		}
		contents.Folders = append(contents.Folders, sf)
	}
	if err := folderRows.Err(); err != nil {
		return nil, err
	}
	return &contents, nil
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
