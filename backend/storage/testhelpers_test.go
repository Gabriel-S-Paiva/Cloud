package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "modernc.org/sqlite"
)

// migrationFilePath resolves the schema file relative to *this source
// file's* location on disk, not the process's working directory. go test
// already runs with cwd = the package directory, so a plain relative path
// would work too — but this is more robust if that ever changes (e.g. a
// test runner that changes cwd, or running a single test file directly).
func migrationFilePath() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "migration", "005_user_root.sql")
}

// newTestStore opens a fresh temp-file SQLite database, applies the full
// schema, and returns a ready-to-use Store. Each call gets its own isolated
// database file that's deleted automatically when the test finishes.
func newTestStore(t *testing.T) *Store {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	schema, err := os.ReadFile(migrationFilePath())
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	return NewStore(db)
}

// seedUser inserts a user directly, bypassing handlers/auth entirely, and
// returns its id. Tests should use this instead of going through the HTTP
// layer — we only want to exercise the storage logic here.
func seedUser(t *testing.T, s *Store, username string, quota int) int {
	t.Helper()

	res, err := s.db.Exec(
		"INSERT INTO Users (username, hashed_password, quota, quota_used) VALUES (?, ?, ?, 0)",
		username, "irrelevant-hash-for-tests", quota,
	)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("seed user id: %v", err)
	}
	return int(id)
}
