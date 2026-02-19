package comments

import "github.com/google/wire"

var (
	CommentsRepoProviders = wire.NewSet(
		ProvideCommentsRepository,
	)
)
