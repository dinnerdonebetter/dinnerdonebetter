package invitations

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideInvitationsService,
		ProvideInvitationDataManager,
		ProvideInvitationDataServer,
	)
)

// ProvideInvitationDataManager turns a database into an InvitationDataManager
func ProvideInvitationDataManager(db database.Database) models.InvitationDataManager {
	return db
}

// ProvideInvitationDataServer is an arbitrary function for dependency injection's sake
func ProvideInvitationDataServer(s *Service) models.InvitationDataServer {
	return s
}
