package handlers

import (
	"backend/storage"
	"context"
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
