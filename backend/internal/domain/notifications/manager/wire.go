package manager

import (
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	notificationsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/notifications"

	"github.com/google/wire"
)

var (
	NotificationsManagerProviders = wire.NewSet(
		NewNotificationsDataManager,
		notificationsrepo.ProvideNotificationsRepository,
		wire.Bind(new(notificationsRepo), new(*notificationsrepo.Repository)),
		wire.Bind(new(notifications.Repository), new(NotificationsDataManager)),
	)
)
