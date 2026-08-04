package handlers

import (
	"backend/middlewares"
	"backend/storage"
	"backend/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type FileHanlder struct {
	store *storage.Store
}
type CreateFileRequest struct {
	DisplayName  string `json:"displayName"`
	ParentFolder *int64 `json:"parentFolder"`
	Size         int    `json:"size"`
	ContentType  string `json:"contentType"`
}
type UpdateFileRequest struct {
	DisplayName  *string `json:"displayName"`
	ParentFolder *int64  `json:"parentFolder"`
}

func NewFileHandler(store *storage.Store) *FileHanlder {
	return &FileHanlder{store: store}
}

func (h *FileHanlder) CreateFile(w http.ResponseWriter, r *http.Request) {
	var newFile CreateFileRequest
	parent := sql.NullInt64{Int64: 0, Valid: false}
	if err := json.NewDecoder(r.Body).Decode(&newFile); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	if newFile.ParentFolder != nil {
		parent = sql.NullInt64{Int64: *newFile.ParentFolder, Valid: true}
	}
	id, err := h.store.CreateFileUploadIntent(r.Context(), newFile.DisplayName, session.UserId, parent, newFile.Size, newFile.ContentType)
	if err != nil {
		utils.WriteJSONError(w, "Error Creating Folder", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (h *FileHanlder) UploadChunk(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "Error parsing id", http.StatusBadRequest)
		return
	}
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	_, permission, err := h.store.FileAccess(r.Context(), id, session.UserId)
	if err != nil {
		utils.WriteJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	if permission != "Owner" {
		utils.WriteJSONError(w, "only the owner can upload to this file", http.StatusForbidden)
		return
	}

	written, err := h.store.AppendFileChunk(r.Context(), id, r.Body)
	if errors.Is(err, storage.ErrFileNotUploading) {
		utils.WriteJSONError(w, "file is not accepting uploads", http.StatusConflict)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"bytesWritten": written})
}

func (h *FileHanlder) GetFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "Error parsing id", http.StatusBadRequest)
		return
	}

	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}

	file, _, err := h.store.FileAccess(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this file", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(file); err != nil {
		utils.WriteJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}

func (h *FileHanlder) GetFileContent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "Error parsing id", http.StatusBadRequest)
		return
	}
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	if _, _, err := h.store.FileAccess(r.Context(), id, session.UserId); err != nil {
		utils.WriteJSONError(w, "You do not own this file", http.StatusForbidden)
		return
	}

	file, content, err := h.store.GetFileContent(r.Context(), id)
	if err != nil {
		utils.WriteJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	defer content.Close()

	w.Header().Set("Content-Type", file.ContentType)
	if r.URL.Query().Get("download") == "true" {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+file.DisplayName+"\"")
	}

	io.Copy(w, content)
}

func (h *FileHanlder) UpdateFile(w http.ResponseWriter, r *http.Request) {
	var update UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		utils.WriteJSONError(w, "Inavlid request body", http.StatusBadRequest)
		return
	}

	var parent *sql.NullInt64
	if update.ParentFolder != nil {
		parent = &sql.NullInt64{Int64: *update.ParentFolder, Valid: true}
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "Error parsing id", http.StatusBadRequest)
		return
	}

	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	_, permission, err := h.store.FileAccess(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this file", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if permission != "Owner" && permission != "Edit" {
		utils.WriteJSONError(w, "you do not have permission to edit this file", http.StatusForbidden)
		return
	}

	err = h.store.UpdateFile(r.Context(), id, update.DisplayName, parent)
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *FileHanlder) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "Error parsing id", http.StatusBadRequest)
		return
	}
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	_, permission, err := h.store.FileAccess(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this file", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if permission != "Owner" && permission != "Edit" {
		utils.WriteJSONError(w, "you do not have permission to delete this file", http.StatusForbidden)
		return
	}

	err = h.store.DeleteFile(r.Context(), id)
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
