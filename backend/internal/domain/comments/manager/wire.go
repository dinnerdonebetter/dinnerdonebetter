package manager

import "github.com/google/wire"

var (
	CommentsManagerProviders = wire.NewSet(
		NewCommentsDataManager,
	)
)
