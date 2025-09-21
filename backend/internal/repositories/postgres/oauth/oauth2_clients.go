package oauth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	types "github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth/generated"
)

const (
	resourceTypeOAuth2Clients = "oauth2_clients"
)

var (
	_ types.OAuth2ClientDataManager = (*repository)(nil)
)

// GetOAuth2ClientByClientID gets an OAuth2 client from the database.
func (r *repository) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*types.OAuth2Client, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if clientID == "" {
		return nil, database.ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientClientIDKey, clientID)
	tracing.AttachToSpan(span, keys.OAuth2ClientClientIDKey, clientID)

	result, err := r.generatedQuerier.GetOAuth2ClientByClientID(ctx, r.db, clientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching oauth2 client")
	}

	client := &types.OAuth2Client{
		CreatedAt:    result.CreatedAt,
		ArchivedAt:   database.TimePointerFromNullTime(result.ArchivedAt),
		Name:         result.Name,
		Description:  result.Description,
		ClientID:     result.ClientID,
		ID:           result.ID,
		ClientSecret: result.ClientSecret,
	}

	return client, nil
}

// GetOAuth2ClientByDatabaseID gets an OAuth2 client from the database.
func (r *repository) GetOAuth2ClientByDatabaseID(ctx context.Context, clientID string) (*types.OAuth2Client, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if clientID == "" {
		return nil, database.ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientClientIDKey, clientID)
	tracing.AttachToSpan(span, keys.OAuth2ClientClientIDKey, clientID)

	result, err := r.generatedQuerier.GetOAuth2ClientByDatabaseID(ctx, r.db, clientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching oauth2 client")
	}

	client := &types.OAuth2Client{
		CreatedAt:    result.CreatedAt,
		ArchivedAt:   database.TimePointerFromNullTime(result.ArchivedAt),
		Name:         result.Name,
		Description:  result.Description,
		ClientID:     result.ClientID,
		ID:           result.ID,
		ClientSecret: result.ClientSecret,
	}

	return client, nil
}

// GetOAuth2Clients gets a list of OAuth2 clients.
func (r *repository) GetOAuth2Clients(ctx context.Context, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.OAuth2Client], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.OAuth2Client]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetOAuth2Clients(ctx, r.db, &generated.GetOAuth2ClientsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching oauth2 clients")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.OAuth2Client{
			CreatedAt:    result.CreatedAt,
			ArchivedAt:   database.TimePointerFromNullTime(result.ArchivedAt),
			Name:         result.Name,
			Description:  result.Description,
			ClientID:     result.ClientID,
			ID:           result.ID,
			ClientSecret: result.ClientSecret,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateOAuth2Client creates an OAuth2 client.
func (r *repository) CreateOAuth2Client(ctx context.Context, input *types.OAuth2ClientDatabaseCreationInput) (*types.OAuth2Client, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}

	logger := r.logger.WithValues(map[string]any{
		keys.OAuth2ClientClientIDKey: input.ClientID,
	})

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if writeErr := r.generatedQuerier.CreateOAuth2Client(ctx, tx, &generated.CreateOAuth2ClientParams{
		ID:           input.ID,
		Name:         input.Name,
		ClientID:     input.ClientID,
		ClientSecret: input.ClientSecret,
	}); writeErr != nil {
		return nil, observability.PrepareError(writeErr, span, "creating OAuth2 client")
	}

	tracing.AttachToSpan(span, keys.OAuth2ClientClientIDKey, input.ID)

	if _, err = r.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: nil,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeOAuth2Clients,
		EventType:        audit.AuditLogEventTypeCreated,
		BelongsToUser:    "",
	}); err != nil {
		r.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	client := &types.OAuth2Client{
		ID:           input.ID,
		Name:         input.Name,
		ClientID:     input.ClientID,
		ClientSecret: input.ClientSecret,
		CreatedAt:    r.CurrentTime(),
	}

	logger.Info("OAuth2 client created")

	return client, nil
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (r *repository) ArchiveOAuth2Client(ctx context.Context, clientID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if clientID == "" {
		return database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.OAuth2ClientClientIDKey, clientID)
	logger := r.logger.WithValue(keys.OAuth2ClientIDKey, clientID)

	rowsAffected, err := r.generatedQuerier.ArchiveOAuth2Client(ctx, r.db, clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return observability.PrepareAndLogError(err, logger, span, "archiving OAuth2 client")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
