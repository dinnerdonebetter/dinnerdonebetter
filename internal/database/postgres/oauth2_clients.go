package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/database"
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

//go:embed queries/oauth2_clients/get_by_client_id.sql
var getOAuth2ClientByClientIDQuery string

// GetOAuth2ClientByClientID gets an OAuth2 client from the database.
func (q *Querier) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*types.OAuth2Client, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if clientID == "" {
		return nil, ErrEmptyInputProvided
	}

	tracing.AttachOAuth2ClientIDToSpan(span, clientID)

	args := []any{clientID}

	row := q.getOneRow(ctx, q.db, "OAuth2 client", getOAuth2ClientByClientIDQuery, args)

	client, _, _, err := q.scanOAuth2Client(ctx, row, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, span, "querying for OAuth2 client")
	}

	return client, nil
}

//go:embed queries/oauth2_clients/get_by_database_id.sql
var getOAuth2ClientByDatabaseIDQuery string

// GetOAuth2ClientByDatabaseID gets an OAuth2 client from the database.
func (q *Querier) GetOAuth2ClientByDatabaseID(ctx context.Context, clientID string) (*types.OAuth2Client, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if clientID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachOAuth2ClientIDToSpan(span, clientID)

	args := []any{clientID}

	row := q.getOneRow(ctx, q.db, "OAuth2 client", getOAuth2ClientByDatabaseIDQuery, args)

	client, _, _, err := q.scanOAuth2Client(ctx, row, false)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, span, "querying for OAuth2 client")
	}

	return client, nil
}

// GetOAuth2Clients gets a list of OAuth2 clients.
func (q *Querier) GetOAuth2Clients(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.OAuth2Client], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)
	x = &types.QueryFilteredResult[types.OAuth2Client]{}
	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "oauth2_clients", nil, nil, nil, userOwnershipColumn, oauth2ClientsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "OAuth2 clients", query, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, observability.PrepareError(err, span, "querying for OAuth2 clients")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanOAuth2Clients(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, span, "scanning response from database")
	}

	return x, nil
}

//go:embed queries/oauth2_clients/create.sql
var createOAuth2ClientQuery string

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

	args := []any{
		input.ID,
		input.Name,
		input.ClientID,
		input.ClientSecret,
	}

	if writeErr := q.performWriteQuery(ctx, q.db, "OAuth2 client creation", createOAuth2ClientQuery, args); writeErr != nil {
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

//go:embed queries/oauth2_clients/archive.sql
var archiveOAuth2ClientQuery string

// ArchiveOAuth2Client archives an OAuth2 client.
func (q *Querier) ArchiveOAuth2Client(ctx context.Context, clientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if clientID == "" {
		return ErrNilInputProvided
	}

	tracing.AttachOAuth2ClientIDToSpan(span, clientID)

	logger := q.logger.WithValues(map[string]any{
		keys.OAuth2ClientIDKey: clientID,
	})

	args := []any{clientID}

	if err := q.performWriteQuery(ctx, q.db, "OAuth2 client archive", archiveOAuth2ClientQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving OAuth2 client")
	}

	logger.Info("OAuth2 client archived")

	return nil
}
