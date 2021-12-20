package service

import (
	"context"
)

// Storager is an interface to interact with the file storage.
type Storager interface {
	// GetFile is used to take file from the storage.
	GetFile(ctx context.Context, path string) ([]byte, error)

	// UploadFile is used to add the file to the storage.
	UploadFile(ctx context.Context, userID int, fileName string, data []byte) (string, error)

	// DeleteFile is used to delete file from the storage.
	DeleteFile(ctx context.Context, path string) error
}
