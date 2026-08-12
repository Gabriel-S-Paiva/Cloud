package storage

import (
	"context"
	"testing"
)

func seedUserWithRoot(t *testing.T, s *Store, username string, quota int) int {
	t.Helper()

	userID := seedUser(t, s, username, quota)

	var rootFolderID int64
	err := s.db.QueryRow(
		"INSERT INTO Folders (display_name, owned_by) VALUES ('root', ?) RETURNING id",
		userID,
	).Scan(&rootFolderID)
	if err != nil {
		t.Fatalf("seed root folder: %v", err)
	}

	if _, err := s.db.Exec("UPDATE Users SET root_folder = ? WHERE id = ?", rootFolderID, userID); err != nil {
		t.Fatalf("attach root folder: %v", err)
	}

	return userID
}

func TestGetUserByID_Success(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUserWithRoot(t, s, "getme", 500)

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if user.Username != "getme" {
		t.Errorf("Username = %q, want %q", user.Username, "getme")
	}
	if user.Quota != 500 {
		t.Errorf("Quota = %d, want 500", user.Quota)
	}
	if user.RootFolderId == 0 {
		t.Errorf("RootFolderId = 0, want a real folder id")
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetUserByID(ctx, 9999)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestCreateRequest_DuplicateUsernameFails(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	if _, err := s.CreateRequest(ctx, "dupeuser", "hash1"); err != nil {
		t.Fatalf("first CreateRequest: %v", err)
	}
	if _, err := s.CreateRequest(ctx, "dupeuser", "hash2"); err != ErrUsernameTaken {
		t.Errorf("expected ErrUsernameTaken, got %v", err)
	}
}

func TestAproveRequest_CreatesUserWithRootFolder(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	requestID, err := s.CreateRequest(ctx, "newbie", "hashedpw")
	if err != nil {
		t.Fatalf("CreateRequest: %v", err)
	}

	userID, err := s.AproveRequest(ctx, requestID)
	if err != nil {
		t.Fatalf("AproveRequest: %v", err)
	}

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if user.Username != "newbie" {
		t.Errorf("Username = %q, want %q", user.Username, "newbie")
	}
	if user.RootFolderId == 0 {
		t.Errorf("RootFolderId = 0, want a real folder id after approval")
	}

	root, err := s.GetFolderById(ctx, user.RootFolderId)
	if err != nil {
		t.Fatalf("GetFolderById(root): %v", err)
	}
	if root.OwnedBy != userID {
		t.Errorf("root folder OwnedBy = %d, want %d", root.OwnedBy, userID)
	}
}

func TestAproveRequest_RemovesRequestRow(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	requestID, err := s.CreateRequest(ctx, "onetime", "hashedpw")
	if err != nil {
		t.Fatalf("CreateRequest: %v", err)
	}
	if _, err := s.AproveRequest(ctx, requestID); err != nil {
		t.Fatalf("AproveRequest: %v", err)
	}

	pending, err := s.ListPendingRequests(ctx)
	if err != nil {
		t.Fatalf("ListPendingRequests: %v", err)
	}
	for _, req := range pending {
		if req.Id == requestID {
			t.Errorf("request %d still pending after approval", requestID)
		}
	}
}

func TestAproveRequest_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.AproveRequest(ctx, 9999)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestRejectRequest_MarksStatusRejected(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	requestID, err := s.CreateRequest(ctx, "rejectme", "hashedpw")
	if err != nil {
		t.Fatalf("CreateRequest: %v", err)
	}

	if err := s.RejectRequest(ctx, requestID); err != nil {
		t.Fatalf("RejectRequest: %v", err)
	}

	pending, err := s.ListPendingRequests(ctx)
	if err != nil {
		t.Fatalf("ListPendingRequests: %v", err)
	}
	for _, req := range pending {
		if req.Id == requestID {
			t.Errorf("rejected request %d still shows as pending", requestID)
		}
	}
}

func TestRejectRequest_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.RejectRequest(ctx, 9999)
	if err != ErrUpdatingRequest {
		t.Errorf("expected ErrUpdatingRequest, got %v", err)
	}
}

func TestListPendingRequests_OnlyReturnsPending(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	pendingID, err := s.CreateRequest(ctx, "stillpending", "hash1")
	if err != nil {
		t.Fatalf("CreateRequest(pending): %v", err)
	}
	rejectedID, err := s.CreateRequest(ctx, "gotrejected", "hash2")
	if err != nil {
		t.Fatalf("CreateRequest(rejected): %v", err)
	}
	if err := s.RejectRequest(ctx, rejectedID); err != nil {
		t.Fatalf("RejectRequest: %v", err)
	}

	pending, err := s.ListPendingRequests(ctx)
	if err != nil {
		t.Fatalf("ListPendingRequests: %v", err)
	}

	found := false
	for _, req := range pending {
		if req.Id == rejectedID {
			t.Errorf("rejected request %d appeared in pending list", rejectedID)
		}
		if req.Id == pendingID {
			found = true
		}
	}
	if !found {
		t.Errorf("pending request %d missing from pending list", pendingID)
	}
}

func TestGetSharableUsers_ExcludesSelf(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	selfID := seedUser(t, s, "self", 1000)
	seedUser(t, s, "other-a", 1000)
	seedUser(t, s, "other-b", 1000)

	sharable, err := s.GetSharableUsers(ctx, selfID)
	if err != nil {
		t.Fatalf("GetSharableUsers: %v", err)
	}
	if len(sharable) != 2 {
		t.Fatalf("got %d sharable users, want 2", len(sharable))
	}
	for _, u := range sharable {
		if u.Id == selfID {
			t.Errorf("GetSharableUsers included the requesting user themselves")
		}
	}
}

func TestSeedAdmin_CreatesAdminWithRootFolder(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	if err := s.SeedAdmin(ctx, "root-admin", "supersecret"); err != nil {
		t.Fatalf("SeedAdmin: %v", err)
	}

	users, err := s.GetUsers(ctx)
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("got %d users after seeding admin, want 1", len(users))
	}
	admin := users[0]
	if admin.Role != "Admin" {
		t.Errorf("Role = %q, want %q", admin.Role, "Admin")
	}
	if admin.RootFolderId == 0 {
		t.Errorf("RootFolderId = 0, want a real folder id")
	}
}

func TestSeedAdmin_IsIdempotent(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	if err := s.SeedAdmin(ctx, "repeat-admin", "pass1"); err != nil {
		t.Fatalf("first SeedAdmin: %v", err)
	}
	if err := s.SeedAdmin(ctx, "repeat-admin", "pass2"); err != nil {
		t.Fatalf("second SeedAdmin: %v", err)
	}

	users, err := s.GetUsers(ctx)
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("got %d users after re-running SeedAdmin, want 1", len(users))
	}
}
