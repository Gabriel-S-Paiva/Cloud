package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"backend/storage"

	_ "modernc.org/sqlite"
)

func migrationFilePath() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "..", "migration", "005_user_root.sql")
}

func newTestServer(t *testing.T) (http.Handler, *storage.Store) {
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

	store := storage.NewStore(db)
	return NewRouter(store), store
}

func doJSON(t *testing.T, handler http.Handler, method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func loginAs(t *testing.T, handler http.Handler, username, password string) *http.Cookie {
	t.Helper()

	body := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	rec := doJSON(t, handler, http.MethodPost, "/login", body, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("login %q: status = %d, body = %s", username, rec.Code, rec.Body.String())
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == "session_token" {
			return c
		}
	}
	t.Fatalf("login %q: response had no session_token cookie", username)
	return nil
}

func seedAdminAndLogin(t *testing.T, store *storage.Store, handler http.Handler) *http.Cookie {
	t.Helper()

	if err := store.SeedAdmin(context.Background(), "admin", "admin-password"); err != nil {
		t.Fatalf("SeedAdmin: %v", err)
	}
	return loginAs(t, handler, "admin", "admin-password")
}

func registerApproveLogin(t *testing.T, handler http.Handler, adminCookie *http.Cookie, username, password string) *http.Cookie {
	t.Helper()

	regBody := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	regRec := doJSON(t, handler, http.MethodPost, "/users", regBody, nil)
	if regRec.Code != http.StatusCreated {
		t.Fatalf("register %q: status = %d, body = %s", username, regRec.Code, regRec.Body.String())
	}
	var regResp struct {
		Id int `json:"id"`
	}
	if err := json.NewDecoder(regRec.Body).Decode(&regResp); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	approvePath := fmt.Sprintf("/users/requests/%d/aprove", regResp.Id)
	approveRec := doJSON(t, handler, http.MethodPost, approvePath, "", adminCookie)
	if approveRec.Code != http.StatusCreated {
		t.Fatalf("approve %q: status = %d, body = %s", username, approveRec.Code, approveRec.Body.String())
	}

	return loginAs(t, handler, username, password)
}

func TestIntegration_RegisterApproveLoginFlow(t *testing.T) {
	handler, store := newTestServer(t)
	adminCookie := seedAdminAndLogin(t, store, handler)

	aliceCookie := registerApproveLogin(t, handler, adminCookie, "alice", "correct horse battery staple")

	rec := doJSON(t, handler, http.MethodGet, "/users/me", "", aliceCookie)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /users/me: status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var me struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&me); err != nil {
		t.Fatalf("decode /users/me: %v", err)
	}
	if me.Username != "alice" {
		t.Errorf("username = %q, want %q", me.Username, "alice")
	}
	if me.Role != "User" {
		t.Errorf("role = %q, want %q", me.Role, "User")
	}
}

func TestIntegration_UnauthenticatedRequestIsRejected(t *testing.T) {
	handler, _ := newTestServer(t)

	rec := doJSON(t, handler, http.MethodGet, "/users/me", "", nil)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("GET /users/me without cookie: status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestIntegration_NonAdminCannotListUsers(t *testing.T) {
	handler, store := newTestServer(t)
	adminCookie := seedAdminAndLogin(t, store, handler)
	aliceCookie := registerApproveLogin(t, handler, adminCookie, "alice", "password123")

	rec := doJSON(t, handler, http.MethodGet, "/users", "", aliceCookie)
	if rec.Code != http.StatusForbidden {
		t.Errorf("GET /users as regular user: status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestIntegration_FolderOwnershipIsEnforced(t *testing.T) {
	handler, store := newTestServer(t)
	adminCookie := seedAdminAndLogin(t, store, handler)
	aliceCookie := registerApproveLogin(t, handler, adminCookie, "alice", "password123")
	bobCookie := registerApproveLogin(t, handler, adminCookie, "bob", "password123")

	createRec := doJSON(t, handler, http.MethodPost, "/folders", `{"displayName":"Alice's Diary"}`, aliceCookie)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create folder: status = %d, body = %s", createRec.Code, createRec.Body.String())
	}
	var created struct {
		Id int64 `json:"id"`
	}
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create folder response: %v", err)
	}
	folderPath := fmt.Sprintf("/folders/%d", created.Id)

	// The owner can read it.
	ownRec := doJSON(t, handler, http.MethodGet, folderPath, "", aliceCookie)
	if ownRec.Code != http.StatusOK {
		t.Errorf("owner GET folder: status = %d, want %d", ownRec.Code, http.StatusOK)
	}

	// A different, unrelated user cannot.
	strangerRec := doJSON(t, handler, http.MethodGet, folderPath, "", bobCookie)
	if strangerRec.Code != http.StatusForbidden {
		t.Errorf("stranger GET folder: status = %d, want %d", strangerRec.Code, http.StatusForbidden)
	}

	// Nobody can without a session at all.
	noAuthRec := doJSON(t, handler, http.MethodGet, folderPath, "", nil)
	if noAuthRec.Code != http.StatusUnauthorized {
		t.Errorf("unauthenticated GET folder: status = %d, want %d", noAuthRec.Code, http.StatusUnauthorized)
	}
}

func TestIntegration_SharedFolderGrantsAccessButNotOwnership(t *testing.T) {
	handler, store := newTestServer(t)
	adminCookie := seedAdminAndLogin(t, store, handler)
	aliceCookie := registerApproveLogin(t, handler, adminCookie, "alice", "password123")
	bobCookie := registerApproveLogin(t, handler, adminCookie, "bob", "password123")

	createRec := doJSON(t, handler, http.MethodPost, "/folders", `{"displayName":"Shared Folder"}`, aliceCookie)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create folder: status = %d, body = %s", createRec.Code, createRec.Body.String())
	}
	var created struct {
		Id int64 `json:"id"`
	}
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create folder response: %v", err)
	}

	summaryRec := doJSON(t, handler, http.MethodGet, "/users/summary", "", aliceCookie)
	if summaryRec.Code != http.StatusOK {
		t.Fatalf("GET /users/summary: status = %d, body = %s", summaryRec.Code, summaryRec.Body.String())
	}
	var sharable []struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(summaryRec.Body).Decode(&sharable); err != nil {
		t.Fatalf("decode /users/summary: %v", err)
	}
	var bobID int
	for _, u := range sharable {
		if u.Username == "bob" {
			bobID = u.Id
		}
	}
	if bobID == 0 {
		t.Fatalf("bob not found in sharable users list: %+v", sharable)
	}

	shareBody := fmt.Sprintf(`{"FolderId":%d,"SharedWith":%d,"Permission":"View"}`, created.Id, bobID)
	shareRec := doJSON(t, handler, http.MethodPost, "/shares", shareBody, aliceCookie)
	if shareRec.Code != http.StatusCreated {
		t.Fatalf("create share: status = %d, body = %s", shareRec.Code, shareRec.Body.String())
	}

	folderPath := fmt.Sprintf("/folders/%d", created.Id)

	bobRead := doJSON(t, handler, http.MethodGet, folderPath, "", bobCookie)
	if bobRead.Code != http.StatusOK {
		t.Errorf("bob GET shared folder: status = %d, want %d", bobRead.Code, http.StatusOK)
	}

	bobDelete := doJSON(t, handler, http.MethodDelete, folderPath, "", bobCookie)
	if bobDelete.Code != http.StatusForbidden {
		t.Errorf("bob DELETE folder (View permission): status = %d, want %d", bobDelete.Code, http.StatusForbidden)
	}
}

func TestIntegration_LoginIsRateLimited(t *testing.T) {
	handler, store := newTestServer(t)
	if err := store.SeedAdmin(context.Background(), "throttled-admin", "correct-password"); err != nil {
		t.Fatalf("SeedAdmin: %v", err)
	}

	badBody := `{"username":"throttled-admin","password":"wrong-password"}`

	var lastCode int
	for i := 0; i < 6; i++ {
		rec := doJSON(t, handler, http.MethodPost, "/login", badBody, nil)
		lastCode = rec.Code
	}

	if lastCode != http.StatusTooManyRequests {
		t.Errorf("6th login attempt within a minute: status = %d, want %d", lastCode, http.StatusTooManyRequests)
	}
}
