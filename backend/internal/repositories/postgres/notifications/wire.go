package notifications

import "github.com/google/wire"

var (
	NotifsRepoProviders = wire.NewSet(
		ProvideNotificationsRepository,
	)
)
