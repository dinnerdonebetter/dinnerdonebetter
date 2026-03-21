package manager

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/random"
)

// RegisterIdentityDataManager registers the identity data manager with the injector.
func RegisterIdentityDataManager(i do.Injector) {
	do.Provide[IdentityDataManager](i, func(i do.Injector) (IdentityDataManager, error) {
		return NewIdentityDataManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[identity.Repository](i),
			do.MustInvoke[random.Generator](i),
			do.MustInvoke[authentication.Hasher](i),
			do.MustInvoke[identityindexing.UserTextSearcher](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
		)
	})
}
