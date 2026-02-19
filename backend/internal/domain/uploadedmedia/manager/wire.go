package manager

import "github.com/google/wire"

var (
	UploadedMediaManagerProviders = wire.NewSet(
		NewUploadedMediaDataManager,
	)
)
