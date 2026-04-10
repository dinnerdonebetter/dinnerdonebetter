package manager

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity"
	identityindexing "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/identity/indexing"

	"github.com/verygoodsoftwarenotvirus/platform/v5/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v5/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v5/random"

	"github.com/samber/do/v2"
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
