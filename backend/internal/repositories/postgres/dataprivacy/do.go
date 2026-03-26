package dataprivacy

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	domaindataprivacy "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/dataprivacy"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/issuereports"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/uploadedmedia"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/waitlists"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/webhooks"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/database"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

// RegisterDataPrivacyRepository registers the data privacy repository with the injector.
func RegisterDataPrivacyRepository(i do.Injector) {
	do.Provide[domaindataprivacy.Repository](i, func(i do.Injector) (domaindataprivacy.Repository, error) {
		return ProvideDataPrivacyRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[audit.Repository](i),
			do.MustInvoke[identity.Repository](i),
			do.MustInvoke[issuereports.Repository](i),
			do.MustInvoke[mealplanning.Repository](i),
			do.MustInvoke[notifications.Repository](i),
			do.MustInvoke[settings.Repository](i),
			do.MustInvoke[uploadedmedia.Repository](i),
			do.MustInvoke[waitlists.Repository](i),
			do.MustInvoke[webhooks.Repository](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
