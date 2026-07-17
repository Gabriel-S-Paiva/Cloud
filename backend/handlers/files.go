package handlers

import (
	"backend/middlewares"
	"backend/storage"
	"backend/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type FileHanlder struct {
	store *storage.Store
}
type CreateFileRequest struct {
	DisplayName  string `json:"displayName"`
	ParentFolder *int64 `json:"parentFolder"`
}
type UpdateFileRequest struct {
	DisplayName  *string `json:"displayName"`
	ParentFolder *int64  `json:"parentFolder"`
}

func NewFileHandler(store *storage.Store) *FileHanlder {
	return &FileHanlder{store: store}
}

func (h *FileHanlder) fileOwnership(ctx context.Context, fileId int, userId int) (*storage.File, error) {
	file, err := h.store.GetFileById(ctx, fileId)
	if err != nil {
		return nil, err
	}
	if file.OwnedBy != userId {
		return nil, storage.ErrForbidden
	}
	return file, nil
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
	_, err := h.store.CreateFile(r.Context(), newFile.DisplayName, session.UserId, parent, 0)
	if err != nil {
		utils.WriteJSONError(w, "Error Creating Folder", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
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

	file, err := h.fileOwnership(r.Context(), id, session.UserId)
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

func (h *FileHanlder) UpdateFile(w http.ResponseWriter, r *http.Request) {
	var update UpdateFileRequest
	parent := sql.NullInt64{Int64: 0, Valid: false}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		utils.WriteJSONError(w, "Inavlid request body", http.StatusBadRequest)
		return
	}
	if update.ParentFolder != nil {
		parent = sql.NullInt64{Int64: *update.ParentFolder, Valid: true}
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
	_, err = h.fileOwnership(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this file", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = h.store.UpdateFile(r.Context(), id, update.DisplayName, &parent)
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
	_, err = h.fileOwnership(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this file", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = h.store.DeleteFile(r.Context(), id)
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
