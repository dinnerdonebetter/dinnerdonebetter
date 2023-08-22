package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.OAuth2ClientDataManager = (*Querier)(nil)

	oauth2ClientsTableColumns = []string{
		"oauth2_clients.id",
		"oauth2_clients.name",
		"oauth2_clients.client_id",
		"oauth2_clients.client_secret",
		"oauth2_clients.created_at",
		"oauth2_clients.archived_at",
	}
)

// scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct.
func (q *Querier) scanOAuth2Client(ctx context.Context, scan database.Scanner, includeCounts bool) (client *types.OAuth2Client, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	client = &types.OAuth2Client{}

	targetVars := []any{
		&client.ID,
		&client.Name,
		&client.ClientID,
		&client.ClientSecret,
		&client.CreatedAt,
		&client.ArchivedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "scanning OAuth2 client database result")
	}

	return client, filteredCount, totalCount, nil
}

// scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2 Clients.
func (q *Querier) scanOAuth2Clients(ctx context.Context, rows database.ResultIterator, includeCounts bool) (clients []*types.OAuth2Client, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		client, fc, tc, scanErr := q.scanOAuth2Client(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, observability.PrepareError(scanErr, span, "scanning OAuth2 client")
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		clients = append(clients, client)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return clients, filteredCount, totalCount, nil
}

// GetOAuth2ClientByClientID gets an OAuth2 client from the database.
func (q *Querier) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*types.OAuth2Client, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if clientID == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientClientIDKey, clientID)
	tracing.AttachOAuth2ClientClientIDToSpan(span, clientID)

	result, err := q.generatedQuerier.GetOAuth2ClientByClientID(ctx, q.db, clientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching oauth2 client")
	}

	client := &types.OAuth2Client{
		CreatedAt:    result.CreatedAt,
		ArchivedAt:   timePointerFromNullTime(result.ArchivedAt),
		Name:         result.Name,
		Description:  result.Description,
		ClientID:     result.ClientID,
		ID:           result.ID,
		ClientSecret: result.ClientSecret,
	}

	return client, nil
}

// GetOAuth2ClientByDatabaseID gets an OAuth2 client from the database.
func (q *Querier) GetOAuth2ClientByDatabaseID(ctx context.Context, clientID string) (*types.OAuth2Client, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if clientID == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientClientIDKey, clientID)
	tracing.AttachOAuth2ClientIDToSpan(span, clientID)

	result, err := q.generatedQuerier.GetOAuth2ClientByDatabaseID(ctx, q.db, clientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching oauth2 client")
	}

	client := &types.OAuth2Client{
		CreatedAt:    result.CreatedAt,
		ArchivedAt:   timePointerFromNullTime(result.ArchivedAt),
		Name:         result.Name,
		Description:  result.Description,
		ClientID:     result.ClientID,
		ID:           result.ID,
		ClientSecret: result.ClientSecret,
	}

	return client, nil
}

// GetOAuth2Clients gets a list of OAuth2 clients.
func (q *Querier) GetOAuth2Clients(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.OAuth2Client], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.OAuth2Client]{
		Pagination: filter.ToPagination(),
	}

	query, args := q.buildListQuery(ctx, "oauth2_clients", nil, nil, nil, userOwnershipColumn, oauth2ClientsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "OAuth2 clients", query, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareAndLogError(err, logger, span, "querying for OAuth2 clients")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanOAuth2Clients(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "scanning response from database")
	}

	return x, nil
}

// CreateOAuth2Client creates an OAuth2 client.
func (q *Querier) CreateOAuth2Client(ctx context.Context, input *types.OAuth2ClientDatabaseCreationInput) (*types.OAuth2Client, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValues(map[string]any{
		keys.OAuth2ClientClientIDKey: input.ClientID,
	})

	if writeErr := q.generatedQuerier.CreateOAuth2Client(ctx, q.db, &generated.CreateOAuth2ClientParams{
		ID:           input.ID,
		Name:         input.Name,
		ClientID:     input.ClientID,
		ClientSecret: input.ClientSecret,
	}); writeErr != nil {
		return nil, observability.PrepareError(writeErr, span, "creating OAuth2 client")
	}

	tracing.AttachOAuth2ClientIDToSpan(span, input.ID)

	client := &types.OAuth2Client{
		ID:           input.ID,
		Name:         input.Name,
		ClientID:     input.ClientID,
		ClientSecret: input.ClientSecret,
		CreatedAt:    q.currentTime(),
	}

	logger.Info("OAuth2 client created")

	return client, nil
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (q *Querier) ArchiveOAuth2Client(ctx context.Context, clientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if clientID == "" {
		return ErrNilInputProvided
	}
	tracing.AttachOAuth2ClientIDToSpan(span, clientID)
	logger := q.logger.WithValue(keys.OAuth2ClientIDKey, clientID)

	if err := q.generatedQuerier.ArchiveOAuth2Client(ctx, q.db, clientID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return observability.PrepareAndLogError(err, logger, span, "archiving OAuth2 client")
	}

	logger.Info("OAuth2 client archived")

	return nil
}
