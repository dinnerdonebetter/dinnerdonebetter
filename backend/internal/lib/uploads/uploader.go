package uploads

import (
	"context"
	"net/http"
)

type (
	// RootUploadDirectory is a type alias for dependency injection's sake.
	RootUploadDirectory string

	// UploadManager stores data in a given storage provider.
	UploadManager interface {
		SaveFile(ctx context.Context, path string, content []byte) error
		ReadFile(ctx context.Context, path string) ([]byte, error)
		ServeFiles(res http.ResponseWriter, req *http.Request)
	}
)
