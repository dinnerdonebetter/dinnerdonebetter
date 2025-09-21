package dataprivacy

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

const (
	resourceTypeUserDataCollections = "user_data_collections"
)

func (r *repository) FetchUserDataCollection(ctx context.Context, userID string) (*dataprivacy.UserDataCollectionResponse, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.WithSpan(span)
	logger.Info("Fetching user data collection")

	// TODO: implement
	x := &dataprivacy.UserDataCollectionResponse{}

	if _, err := r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.DB(), &audit.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: nil,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeUserDataCollections,
		EventType:        audit.AuditLogEventTypeCreated,
		BelongsToUser:    "",
	}); err != nil {
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	return x, nil
}

func (r *repository) DeleteUser(ctx context.Context, userID string) error {
	// TODO: implement
	return nil
}
