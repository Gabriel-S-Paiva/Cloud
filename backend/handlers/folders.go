package handlers

import (
	"backend/middlewares"
	"backend/storage"
	"backend/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type FolderHandler struct {
	store *storage.Store
}
type CreateFolderRequest struct {
	DisplayName  string `json:"displayName"`
	ParentFolder *int64 `json:"parentFolder"`
}
type UpdateFolderRequest struct {
	DisplayName  *string `json:"displayName"`
	ParentFolder *int64  `json:"parentFolder"`
}

func NewFolderHandlers(store *storage.Store) *FolderHandler {
	return &FolderHandler{store: store}
}

func (h *FolderHandler) folderOwnership(ctx context.Context, folderId int, userId int) (*storage.Folder, error) {
	folder, err := h.store.GetFolderById(ctx, folderId)
	if err != nil {
		return nil, err
	}
	if folder.OwnedBy != userId {
		return nil, storage.ErrForbidden
	}
	return folder, nil
}

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	var newFolder CreateFolderRequest
	parent := sql.NullInt64{Int64: 0, Valid: false}
	if err := json.NewDecoder(r.Body).Decode(&newFolder); err != nil {
		utils.WriteJSONError(w, "Inavlid request body", http.StatusBadRequest)
		return
	}
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	if newFolder.ParentFolder != nil {
		parent = sql.NullInt64{Int64: *newFolder.ParentFolder, Valid: true}
	}
	_, err := h.store.CreateFolder(r.Context(), newFolder.DisplayName, session.UserId, parent)
	if err != nil {
		utils.WriteJSONError(w, "Error Creating Folder", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *FolderHandler) GetFolder(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.folderOwnership(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this folder", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	folder, err := h.store.GetFolderById(r.Context(), id)
	if err != nil {
		utils.WriteJSONError(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	if err := json.NewEncoder(w).Encode(folder); err != nil {
		utils.WriteJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}

func (h *FolderHandler) GetFolderContents(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.folderOwnership(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this folder", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	folder, err := h.store.GetFolderContents(r.Context(), id)
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(folder); err != nil {
		utils.WriteJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}

func (h *FolderHandler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	var update UpdateFolderRequest
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
	_, err = h.folderOwnership(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this folder", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = h.store.UpdateFolder(r.Context(), id, update.DisplayName, &parent)
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.folderOwnership(r.Context(), id, session.UserId)
	if errors.Is(err, storage.ErrForbidden) {
		utils.WriteJSONError(w, "You do not own this folder", http.StatusForbidden)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = h.store.DeleteFolder(r.Context(), id)
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
