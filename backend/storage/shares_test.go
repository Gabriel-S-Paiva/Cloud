package storage

import (
	"context"
	"database/sql"
	"testing"
)

func TestCreateShare_ThenGetShareById(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "shareowner1", 1000)
	targetID := seedUser(t, s, "sharetarget1", 1000)

	folderID, err := s.CreateFolder(ctx, "ToShare", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}

	shareID, err := s.CreateShare(ctx, sql.NullInt64{}, sql.NullInt64{Int64: folderID, Valid: true}, targetID, "View")
	if err != nil {
		t.Fatalf("CreateShare: %v", err)
	}

	share, err := s.GetShareById(ctx, shareID)
	if err != nil {
		t.Fatalf("GetShareById: %v", err)
	}
	if share.SharedWith != targetID {
		t.Errorf("SharedWith = %d, want %d", share.SharedWith, targetID)
	}
	if share.Permission != "View" {
		t.Errorf("Permission = %q, want %q", share.Permission, "View")
	}
	if !share.Folder.Valid || share.Folder.Int64 != folderID {
		t.Errorf("Folder = %+v, want valid %d", share.Folder, folderID)
	}
	if share.File.Valid {
		t.Errorf("File.Valid = true, want false (this is a folder share)")
	}
}

func TestGetShareById_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetShareById(ctx, 9999)
	if err != ErrShareNotFound {
		t.Errorf("expected ErrShareNotFound, got %v", err)
	}
}

func TestGetIncomingShares_ReturnsSharesForRecipient(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "incowner", 10000)
	recipientID := seedUser(t, s, "increcipient", 1000)

	folderID, err := s.CreateFolder(ctx, "SharedFolder", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}
	fileID, err := s.CreateFileUploadIntent(ctx, "shared.txt", ownerID, sql.NullInt64{}, 5, "text/plain")
	if err != nil {
		t.Fatalf("CreateFileUploadIntent: %v", err)
	}

	if _, err := s.CreateShare(ctx, sql.NullInt64{}, sql.NullInt64{Int64: folderID, Valid: true}, recipientID, "View"); err != nil {
		t.Fatalf("CreateShare(folder): %v", err)
	}
	if _, err := s.CreateShare(ctx, sql.NullInt64{Int64: fileID, Valid: true}, sql.NullInt64{}, recipientID, "Edit"); err != nil {
		t.Fatalf("CreateShare(file): %v", err)
	}

	incoming, err := s.GetIncomingShares(ctx, recipientID)
	if err != nil {
		t.Fatalf("GetIncomingShares: %v", err)
	}
	if len(incoming.Folders) != 1 {
		t.Errorf("got %d shared folders, want 1", len(incoming.Folders))
	}
	if len(incoming.Files) != 1 {
		t.Errorf("got %d shared files, want 1", len(incoming.Files))
	}
	if len(incoming.Files) == 1 && incoming.Files[0].OwnedByUsername != "incowner" {
		t.Errorf("OwnedByUsername = %q, want %q", incoming.Files[0].OwnedByUsername, "incowner")
	}
}

func TestGetOutgoingShares_ReturnsSharesCreatedByOwner(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "outowner", 1000)
	recipientID := seedUser(t, s, "outrecipient", 1000)

	folderID, err := s.CreateFolder(ctx, "MyShare", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}
	if _, err := s.CreateShare(ctx, sql.NullInt64{}, sql.NullInt64{Int64: folderID, Valid: true}, recipientID, "View"); err != nil {
		t.Fatalf("CreateShare: %v", err)
	}

	outgoing, err := s.GetOutgoingShares(ctx, ownerID)
	if err != nil {
		t.Fatalf("GetOutgoingShares: %v", err)
	}
	if len(outgoing.Folders) != 1 {
		t.Fatalf("got %d outgoing folder shares, want 1", len(outgoing.Folders))
	}
	if outgoing.Folders[0].SharedWith != "outrecipient" {
		t.Errorf("SharedWith = %q, want %q", outgoing.Folders[0].SharedWith, "outrecipient")
	}

	// Sanity check: the recipient shouldn't see this in their own outgoing list.
	recipientOutgoing, err := s.GetOutgoingShares(ctx, recipientID)
	if err != nil {
		t.Fatalf("GetOutgoingShares(recipient): %v", err)
	}
	if len(recipientOutgoing.Folders) != 0 {
		t.Errorf("recipient's outgoing folders = %d, want 0", len(recipientOutgoing.Folders))
	}
}

func TestUpdatePermission_ChangesShareLevel(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "permowner", 1000)
	targetID := seedUser(t, s, "permtarget", 1000)

	folderID, err := s.CreateFolder(ctx, "PermFolder", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}
	shareID, err := s.CreateShare(ctx, sql.NullInt64{}, sql.NullInt64{Int64: folderID, Valid: true}, targetID, "View")
	if err != nil {
		t.Fatalf("CreateShare: %v", err)
	}

	if err := s.UpdatePermission(ctx, shareID, "Edit"); err != nil {
		t.Fatalf("UpdatePermission: %v", err)
	}

	updated, err := s.GetShareById(ctx, shareID)
	if err != nil {
		t.Fatalf("GetShareById: %v", err)
	}
	if updated.Permission != "Edit" {
		t.Errorf("Permission = %q, want %q", updated.Permission, "Edit")
	}
}

func TestUpdatePermission_NonexistentIdReturnsError(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.UpdatePermission(ctx, 9999, "Edit")
	if err != ErrUpdatingShare {
		t.Errorf("expected ErrUpdatingShare, got %v", err)
	}
}

func TestDeleteShare_Success(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	ownerID := seedUser(t, s, "delshareowner", 1000)
	targetID := seedUser(t, s, "delsharetarget", 1000)

	folderID, err := s.CreateFolder(ctx, "DelShareFolder", ownerID, sql.NullInt64{})
	if err != nil {
		t.Fatalf("CreateFolder: %v", err)
	}
	shareID, err := s.CreateShare(ctx, sql.NullInt64{}, sql.NullInt64{Int64: folderID, Valid: true}, targetID, "View")
	if err != nil {
		t.Fatalf("CreateShare: %v", err)
	}

	if err := s.DeleteShare(ctx, shareID); err != nil {
		t.Fatalf("DeleteShare: %v", err)
	}
	if _, err := s.GetShareById(ctx, shareID); err != ErrShareNotFound {
		t.Errorf("expected ErrShareNotFound after delete, got %v", err)
	}
}

func TestDeleteShare_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.DeleteShare(ctx, 9999)
	if err != ErrShareNotFound {
		t.Errorf("expected ErrShareNotFound, got %v", err)
	}
}
