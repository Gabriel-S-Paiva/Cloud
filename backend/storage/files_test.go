package storage

import (
	"bytes"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestCreateFileUploadIntent_QuotaEnforcement(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "quotauser", 1000)

	t.Run("within quota succeeds", func(t *testing.T) {
		_, err := s.CreateFileUploadIntent(ctx, "small.txt", userID, sql.NullInt64{}, 500, "text/plain")
		if err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})

	t.Run("exceeding remaining quota fails", func(t *testing.T) {
		// 500 already used above, quota is 1000 — 600 more should not fit.
		_, err := s.CreateFileUploadIntent(ctx, "big.txt", userID, sql.NullInt64{}, 600, "text/plain")
		if err != ErrQuotaExceeded {
			t.Errorf("expected ErrQuotaExceeded, got %v", err)
		}
	})
}

func TestGetFileById_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetFileById(ctx, 9999)
	if err != ErrFileNotFound {
		t.Errorf("expected ErrFileNotFound, got %v", err)
	}
}

func TestAppendFileChunk_MarksCompleteWhenFullyReceived(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "uploaduser", 10000)

	content := []byte("hello world")
	fileID, err := s.CreateFileUploadIntent(ctx, "test.txt", userID, sql.NullInt64{}, len(content), "text/plain")
	if err != nil {
		t.Fatalf("create intent: %v", err)
	}
	t.Cleanup(func() {
		os.Remove(filepath.Join(storageDir, strconv.FormatInt(fileID, 10)))
	})

	written, err := s.AppendFileChunk(ctx, int(fileID), bytes.NewReader(content))
	if err != nil {
		t.Fatalf("append chunk: %v", err)
	}
	if int(written) != len(content) {
		t.Errorf("wrote %d bytes, want %d", written, len(content))
	}

	file, err := s.GetFileById(ctx, int(fileID))
	if err != nil {
		t.Fatalf("get file: %v", err)
	}
	if file.Status != "Complete" {
		t.Errorf("status = %q, want %q", file.Status, "Complete")
	}
}

// TestDeleteFile_RemovesBinaryFromDisk is the regression test for the bug
// where DeleteFile removed the DB row but left the file's bytes on disk.
func TestDeleteFile_RemovesBinaryFromDisk(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "deleteuser", 10000)

	content := []byte("delete me")
	fileID, err := s.CreateFileUploadIntent(ctx, "test.txt", userID, sql.NullInt64{}, len(content), "text/plain")
	if err != nil {
		t.Fatalf("create intent: %v", err)
	}

	if _, err := s.AppendFileChunk(ctx, int(fileID), bytes.NewReader(content)); err != nil {
		t.Fatalf("append chunk: %v", err)
	}

	path := filepath.Join(storageDir, strconv.FormatInt(fileID, 10))
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("setup: file should exist on disk before delete, stat err: %v", err)
	}
	// Belt-and-braces cleanup in case the assertion below fails and the
	// file is left behind.
	t.Cleanup(func() { os.Remove(path) })

	if err := s.DeleteFile(ctx, int(fileID)); err != nil {
		t.Fatalf("DeleteFile: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected file to be removed from disk after DeleteFile, but it still exists")
	}

	if _, err := s.GetFileById(ctx, int(fileID)); err != ErrFileNotFound {
		t.Errorf("expected row to be gone (ErrFileNotFound), got %v", err)
	}
}

func TestDeleteFile_RefundsQuota(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	userID := seedUser(t, s, "refunduser", 1000)

	fileID, err := s.CreateFileUploadIntent(ctx, "test.txt", userID, sql.NullInt64{}, 400, "text/plain")
	if err != nil {
		t.Fatalf("create intent: %v", err)
	}

	var quotaUsedBefore int
	if err := s.db.QueryRow("SELECT quota_used FROM Users WHERE id = ?", userID).Scan(&quotaUsedBefore); err != nil {
		t.Fatalf("read quota_used: %v", err)
	}
	if quotaUsedBefore != 400 {
		t.Fatalf("setup: quota_used = %d, want 400", quotaUsedBefore)
	}

	if err := s.DeleteFile(ctx, int(fileID)); err != nil {
		t.Fatalf("DeleteFile: %v", err)
	}

	var quotaUsedAfter int
	if err := s.db.QueryRow("SELECT quota_used FROM Users WHERE id = ?", userID).Scan(&quotaUsedAfter); err != nil {
		t.Fatalf("read quota_used after delete: %v", err)
	}
	if quotaUsedAfter != 0 {
		t.Errorf("quota_used after delete = %d, want 0", quotaUsedAfter)
	}
}
