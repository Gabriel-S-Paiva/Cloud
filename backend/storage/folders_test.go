package storage

import (
	"context"
	"database/sql"
	"testing"
)

func TestCreateFolder_ThenGetFolderById(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "folderuser", 1000)

	id, err := s.CreateFolder(ctx, "Documents", userID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}

	folder, err := s.GetFolderById(ctx, int(id))
	if err != nil {
		t.Fatalf("GetFolderById: %v", err)
	}
	if folder.DisplayName != "Documents" {
		t.Errorf("DisplayName = %q, want %q", folder.DisplayName, "Documents")
	}
	if folder.OwnedBy != userID {
		t.Errorf("OwnedBy = %d, want %d", folder.OwnedBy, userID)
	}
	if folder.ParentFolder.Valid {
		t.Errorf("ParentFolder.Valid = true, want false (root-level folder)")
	}
}

func TestGetFolderById_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetFolderById(ctx, 9999)
	if err != ErrFolderNotFound {
		t.Errorf("expected ErrFolderNotFound, got %v", err)
	}
}

func TestGetFolderContents_ReturnsFilesAndSubfolders(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "contentsuser", 10000)

	parentID, err := s.CreateFolder(ctx, "Parent", userID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder(parent): %v", err)
	}
	parentRef := sql.NullInt64{Int64: parentID, Valid: true}

	if _, err := s.CreateFolder(ctx, "Child", userID, parentRef); err != nil {
		t.Fatalf("CreateFolder(child): %v", err)
	}
	if _, err := s.CreateFileUploadIntent(ctx, "notes.txt", userID, parentRef, 10, "text/plain"); err != nil {
		t.Fatalf("CreateFileUploadIntent: %v", err)
	}

	contents, err := s.GetFolderContents(ctx, int(parentID))
	if err != nil {
		t.Fatalf("GetFolderContents: %v", err)
	}
	if len(contents.Folders) != 1 {
		t.Errorf("got %d subfolders, want 1", len(contents.Folders))
	}
	if len(contents.Files) != 1 {
		t.Errorf("got %d files, want 1", len(contents.Files))
	}
}

func TestUpdateFolder_RenameAndMove(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "moveuser", 1000)

	destID, err := s.CreateFolder(ctx, "Destination", userID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder(dest): %v", err)
	}
	folderID, err := s.CreateFolder(ctx, "Old Name", userID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder(folder): %v", err)
	}

	newName := "New Name"
	newParent := sql.NullInt64{Int64: destID, Valid: true}
	if err := s.UpdateFolder(ctx, int(folderID), &newName, &newParent); err != nil {
		t.Fatalf("UpdateFolder: %v", err)
	}

	updated, err := s.GetFolderById(ctx, int(folderID))
	if err != nil {
		t.Fatalf("GetFolderById: %v", err)
	}
	if updated.DisplayName != "New Name" {
		t.Errorf("DisplayName = %q, want %q", updated.DisplayName, "New Name")
	}
	if !updated.ParentFolder.Valid || updated.ParentFolder.Int64 != destID {
		t.Errorf("ParentFolder = %+v, want valid %d", updated.ParentFolder, destID)
	}
}

func TestUpdateFolder_NonexistentIdReturnsError(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	newName := "Whatever"
	err := s.UpdateFolder(ctx, 9999, &newName, nil)
	if err != ErrUpdatingFolder {
		t.Errorf("expected ErrUpdatingFolder, got %v", err)
	}
}

func TestDeleteFolder_Success(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "deletefolderuser", 1000)

	id, err := s.CreateFolder(ctx, "ToDelete", userID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}

	if err := s.DeleteFolder(ctx, int(id)); err != nil {
		t.Fatalf("DeleteFolder: %v", err)
	}
	if _, err := s.GetFolderById(ctx, int(id)); err != ErrFolderNotFound {
		t.Errorf("expected ErrFolderNotFound after delete, got %v", err)
	}
}

func TestDeleteFolder_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.DeleteFolder(ctx, 9999)
	if err != ErrFolderNotFound {
		t.Errorf("expected ErrFolderNotFound, got %v", err)
	}
}

func TestFolderOwnership_ForbidsNonOwner(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "owner1", 1000)
	otherID := seedUser(t, s, "other1", 1000)

	folderID, err := s.CreateFolder(ctx, "Private", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}

	if _, err := s.FolderOwnership(ctx, int(folderID), ownerID); err != nil {
		t.Errorf("owner: expected success, got %v", err)
	}
	if _, err := s.FolderOwnership(ctx, int(folderID), otherID); err != ErrForbidden {
		t.Errorf("non-owner: expected ErrForbidden, got %v", err)
	}
}

func TestFolderAccess_OwnerGetsOwnerPermission(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "owner2", 1000)

	folderID, err := s.CreateFolder(ctx, "Mine", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}

	_, permission, err := s.FolderAccess(ctx, int(folderID), ownerID)
	if err != nil {
		t.Fatalf("FolderAccess: %v", err)
	}
	if permission != "Owner" {
		t.Errorf("permission = %q, want %q", permission, "Owner")
	}
}

func TestFolderAccess_SharedUserGetsSharePermission(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "owner3", 1000)
	sharedID := seedUser(t, s, "shared3", 1000)

	folderID, err := s.CreateFolder(ctx, "Shared", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}
	folderRef := sql.NullInt64{Int64: folderID, Valid: true}
	if _, err := s.CreateShare(ctx, sql.NullInt64{}, folderRef, sharedID, "Edit"); err != nil {
		t.Fatalf("CreateShare: %v", err)
	}

	_, permission, err := s.FolderAccess(ctx, int(folderID), sharedID)
	if err != nil {
		t.Fatalf("FolderAccess: %v", err)
	}
	if permission != "Edit" {
		t.Errorf("permission = %q, want %q", permission, "Edit")
	}
}

func TestFolderAccess_NoAccessReturnsError(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "owner4", 1000)
	strangerID := seedUser(t, s, "stranger4", 1000)

	folderID, err := s.CreateFolder(ctx, "NotYours", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}

	if _, _, err := s.FolderAccess(ctx, int(folderID), strangerID); err != ErrFileNotFound {
		t.Errorf("expected ErrFileNotFound, got %v", err)
	}
}
