package storage

import (
	"context"
	"database/sql"
)

type Folder struct {
	Id           int
	DisplayName  string
	OwnedBy      int
	ParentFolder sql.NullInt64
}
type FolderContents struct {
	Folders []Folder
	Files   []File
}

func (s *Store) CreateFolder(ctx context.Context, displayName string, ownedBy int, parentFolder sql.NullInt64) (int64, error) {
	result, err := s.db.ExecContext(ctx, "INSERT INTO Folders (display_name, owned_by, parent_folder) VALUES (?,?,?)", displayName, ownedBy, parentFolder)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) GetFolderById(ctx context.Context, id int) (*Folder, error) {
	var folder Folder
	err := s.db.QueryRowContext(ctx, "SELECT id, display_name, owned_by, parent_folder FROM Folders WHERE id = ?", id).Scan(&folder.Id, &folder.DisplayName, &folder.OwnedBy, &folder.ParentFolder)
	if err == sql.ErrNoRows {
		return nil, ErrFolderNotFound
	}
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (s *Store) GetFolderContents(ctx context.Context, id int) (*FolderContents, error) {
	var folderContents FolderContents
	fileRows, err := s.db.QueryContext(ctx, "SELECT id, display_name, owned_by, size, uploaded_at, last_modified, parent_folder FROM Files WHERE parent_folder = ?", id)
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

	folderRows, err := s.db.QueryContext(ctx, "SELECT id, display_name, owned_by, parent_folder FROM Folders WHERE parent_folder = ?", id)
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

func (s *Store) UpdateFolder(ctx context.Context, id int, newName *string, newParent *sql.NullInt64) error {
	if newParent != nil {
		res, err := s.db.ExecContext(ctx, "UPDATE Folders SET parent_folder = ? WHERE id = ?", newParent, id)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows != 1 {
			return ErrUpdatingFolder
		}
	}
	if newName != nil {
		res, err := s.db.ExecContext(ctx, "UPDATE Folders SET display_name = ? WHERE id = ?", newName, id)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows != 1 {
			return ErrUpdatingFolder
		}
	}
	return nil
}

func (s *Store) DeleteFolder(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM Folders WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrFolderNotFound
	}
	if rows != 1 {
		return ErrFolderMismatch
	}
	return nil
}
