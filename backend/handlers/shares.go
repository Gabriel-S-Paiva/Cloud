package handlers

import (
	"backend/middlewares"
	"backend/storage"
	"backend/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type CreateShareRequest struct {
	FileId     *int64
	FolderId   *int64
	SharedWith int
	Permission string
}

type UpdateShareRequest struct {
	Permission string
}

type ShareHandler struct {
	store *storage.Store
}

func NewShareHandlers(store *storage.Store) *ShareHandler {
	return &ShareHandler{store: store}
}

func (h *ShareHandler) CreateShare(w http.ResponseWriter, r *http.Request) {
	var newShare CreateShareRequest
	file := sql.NullInt64{Int64: 0, Valid: false}
	folder := sql.NullInt64{Int64: 0, Valid: false}

	if err := json.NewDecoder(r.Body).Decode(&newShare); err != nil {
		utils.WriteJSONError(w, "Invalid Request body", http.StatusBadRequest)
		return
	}
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	if (newShare.FileId == nil) == (newShare.FolderId == nil) {
		utils.WriteJSONError(w, "exactly one of fileId or folderId must be provided", http.StatusBadRequest)
		return
	}
	if newShare.FileId != nil {
		if _, err := h.store.FileOwnership(r.Context(), int(*newShare.FileId), session.UserId); err != nil {
			utils.WriteJSONError(w, "You don't own this file", http.StatusForbidden)
			return
		}
		file = sql.NullInt64{Int64: *newShare.FileId, Valid: true}
	}
	if newShare.FolderId != nil {
		if _, err := h.store.FolderOwnership(r.Context(), int(*newShare.FolderId), session.UserId); err != nil {
			utils.WriteJSONError(w, "You don't own this folder", http.StatusForbidden)
			return
		}
		folder = sql.NullInt64{Int64: *newShare.FolderId, Valid: true}
	}
	id, err := h.store.CreateShare(r.Context(), file, folder, newShare.SharedWith, newShare.Permission)
	if err != nil {
		utils.WriteJSONError(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// Shared With Me
func (h *ShareHandler) ViewIncomingShares(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	contents, err := h.store.GetIncomingShares(r.Context(), session.UserId)
	if err != nil {
		utils.WriteJSONError(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contents); err != nil {
		utils.WriteJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}

// Shared Files
func (h *ShareHandler) ViewOutgoingShares(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}
	contents, err := h.store.GetOutgoingShares(r.Context(), session.UserId)
	if err != nil {
		utils.WriteJSONError(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contents); err != nil {
		utils.WriteJSONError(w, "could not encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ShareHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
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
	share, err := h.store.GetShareById(r.Context(), id)
	if errors.Is(err, storage.ErrShareNotFound) {
		utils.WriteJSONError(w, "Share not found", http.StatusNotFound)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Error Fetching Share", http.StatusInternalServerError)
		return
	}
	if share.File.Valid {
		if _, err := h.store.FileOwnership(r.Context(), int(share.File.Int64), session.UserId); err != nil {
			utils.WriteJSONError(w, "you don't own this share", http.StatusForbidden)
			return
		}
	}
	if share.Folder.Valid {
		if _, err := h.store.FolderOwnership(r.Context(), int(share.Folder.Int64), session.UserId); err != nil {
			utils.WriteJSONError(w, "you don't own this share", http.StatusForbidden)
			return
		}
	}
	var permission UpdateShareRequest
	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		utils.WriteJSONError(w, "Invalid Requets Body", http.StatusBadRequest)
		return
	}
	err = h.store.UpdatePermission(r.Context(), id, permission.Permission)
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ShareHandler) DeleteShare(w http.ResponseWriter, r *http.Request) {
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
	share, err := h.store.GetShareById(r.Context(), id)
	if errors.Is(err, storage.ErrShareNotFound) {
		utils.WriteJSONError(w, "Share not found", http.StatusNotFound)
		return
	}
	if err != nil {
		utils.WriteJSONError(w, "Error Fetching Share", http.StatusInternalServerError)
		return
	}

	isRecipient := share.SharedWith == session.UserId
	isOwner := false
	if share.File.Valid {
		if _, err := h.store.FileOwnership(r.Context(), int(share.File.Int64), session.UserId); err == nil {
			isOwner = true
		}
	}
	if share.Folder.Valid {
		if _, err := h.store.FolderOwnership(r.Context(), int(share.Folder.Int64), session.UserId); err == nil {
			isOwner = true
		}
	}

	if !isOwner && !isRecipient {
		utils.WriteJSONError(w, "you don't have access to this share", http.StatusForbidden)
		return
	}

	err = h.store.DeleteShare(r.Context(), id)
	if err != nil {
		utils.WriteJSONError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
