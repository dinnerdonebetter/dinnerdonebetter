package manager

import (
	"github.com/dinnerdonebetter/backend/internal/domain/waitlists"
	waitlistsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/waitlists"

	"github.com/google/wire"
)

var (
	WaitlistManagerProviders = wire.NewSet(
		NewWaitlistDataManager,
		waitlistsrepo.ProvideWaitlistsRepository,
		wire.Bind(new(waitlistRepository), new(*waitlistsrepo.Repository)),
		wire.Bind(new(waitlists.Repository), new(WaitlistsDataManager)),
	)
)
