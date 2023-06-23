package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var _ types.OAuth2ClientTokenDataManager = (*Querier)(nil)

// scanOAuth2ClientToken takes a database Scanner (i.e. *sql.Row) and scans the result into an OAuth2 client token struct.
func (q *Querier) scanOAuth2ClientToken(ctx context.Context, scan database.Scanner) (x *types.OAuth2ClientToken, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.OAuth2ClientToken{}

	targetVars := []any{
		&x.ID,
		&x.ClientID,
		&x.BelongsToUser,
		&x.RedirectURI,
		&x.Scope,
		&x.Code,
		&x.CodeChallenge,
		&x.CodeChallengeMethod,
		&x.CodeCreateAt,
		&x.CodeExpiresIn,
		&x.Access,
		&x.AccessCreateAt,
		&x.AccessExpiresIn,
		&x.Refresh,
		&x.RefreshCreateAt,
		&x.RefreshExpiresIn,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, span, "")
	}

	return x, nil
}

//go:embed queries/oauth2_client_tokens/get_one_by_code.sql
var getOAuth2ClientTokenByCodeQuery string

// GetOAuth2ClientTokenByCode fetches an OAuth2 client token from the database.
func (q *Querier) GetOAuth2ClientTokenByCode(ctx context.Context, code string) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if code == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientTokenCodeKey, code)
	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenCodeKey, code)

	args := []any{
		code,
	}

	row := q.getOneRow(ctx, q.db, "oauth2 client token by code", getOAuth2ClientTokenByCodeQuery, args)

	oauth2ClientToken, err := q.scanOAuth2ClientToken(ctx, row)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning oauth2 client token")
	}

	return oauth2ClientToken, nil
}

//go:embed queries/oauth2_client_tokens/get_one_by_access.sql
var getOAuth2ClientTokenByAccessQuery string

// GetOAuth2ClientTokenByAccess fetches an OAuth2 client token from the database.
func (q *Querier) GetOAuth2ClientTokenByAccess(ctx context.Context, code string) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if code == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientTokenAccessKey, code)
	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenAccessKey, code)

	args := []any{
		code,
	}

	row := q.getOneRow(ctx, q.db, "oauth2 client token by access", getOAuth2ClientTokenByAccessQuery, args)

	oauth2ClientToken, err := q.scanOAuth2ClientToken(ctx, row)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning oauth2 client token")
	}

	return oauth2ClientToken, nil
}

//go:embed queries/oauth2_client_tokens/get_one_by_refresh.sql
var getOAuth2ClientTokenByRefreshQuery string

// GetOAuth2ClientTokenByRefresh fetches an OAuth2 client token from the database.
func (q *Querier) GetOAuth2ClientTokenByRefresh(ctx context.Context, code string) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if code == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientTokenRefreshKey, code)
	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenRefreshKey, code)

	args := []any{
		code,
	}

	row := q.getOneRow(ctx, q.db, "oauth2 client token by refresh", getOAuth2ClientTokenByRefreshQuery, args)

	oauth2ClientToken, err := q.scanOAuth2ClientToken(ctx, row)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning oauth2 client token")
	}

	return oauth2ClientToken, nil
}

//go:embed queries/oauth2_client_tokens/create.sql
var oauth2ClientTokenCreationQuery string

// CreateOAuth2ClientToken creates an OAuth2 client token in the database.
func (q *Querier) CreateOAuth2ClientToken(ctx context.Context, input *types.OAuth2ClientTokenDatabaseCreationInput) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.OAuth2ClientTokenIDKey, input.ID)
	now := q.currentTime()

	args := []any{
		input.ID,
		input.ClientID,
		input.BelongsToUser,
		input.RedirectURI,
		input.Scope,
		input.Code,
		input.CodeChallenge,
		input.CodeChallengeMethod,
		now,
		now.Add(input.CodeExpiresIn),
		input.Access,
		now,
		now.Add(input.AccessExpiresIn),
		input.Refresh,
		now,
		now.Add(input.RefreshExpiresIn),
	}

	// create the oauth2 client token.
	if err := q.performWriteQuery(ctx, q.db, "oauth2 client token creation", oauth2ClientTokenCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing oauth2 client token creation query")
	}

	x := &types.OAuth2ClientToken{
		ID:                  input.ID,
		ClientID:            input.ClientID,
		BelongsToUser:       input.BelongsToUser,
		RedirectURI:         input.RedirectURI,
		Scope:               input.Scope,
		Code:                input.Code,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
		CodeCreateAt:        input.CodeCreateAt,
		CodeExpiresIn:       input.CodeExpiresIn,
		Access:              input.Access,
		AccessCreateAt:      input.AccessCreateAt,
		AccessExpiresIn:     input.AccessExpiresIn,
		Refresh:             input.Refresh,
		RefreshCreateAt:     input.RefreshCreateAt,
		RefreshExpiresIn:    input.RefreshExpiresIn,
	}

	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenIDKey, x.ID)
	logger.Info("oauth2 client token created")

	return x, nil
}

//go:embed queries/oauth2_client_tokens/archive_by_access.sql
var archiveOAuth2ClientTokenByAccessQuery string

// ArchiveOAuth2ClientTokenByAccess archives an OAuth2 client token from the database by its ID.
func (q *Querier) ArchiveOAuth2ClientTokenByAccess(ctx context.Context, access string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if access == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientTokenAccessKey, access)
	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenAccessKey, access)

	args := []any{
		access,
	}

	if err := q.performWriteQuery(ctx, q.db, "oauth2 client token archive by access", archiveOAuth2ClientTokenByAccessQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating oauth2 client token by access")
	}

	logger.Info("oauth2 client token archived by access")

	return nil
}

//go:embed queries/oauth2_client_tokens/archive_by_code.sql
var archiveOAuth2ClientTokenByCodeQuery string

// ArchiveOAuth2ClientTokenByCode archives an OAuth2 client token from the database by its ID.
func (q *Querier) ArchiveOAuth2ClientTokenByCode(ctx context.Context, code string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if code == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientTokenCodeKey, code)
	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenCodeKey, code)

	args := []any{
		code,
	}

	if err := q.performWriteQuery(ctx, q.db, "oauth2 client token archive by code", archiveOAuth2ClientTokenByCodeQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating oauth2 client token by code")
	}

	logger.Info("oauth2 client token archived by code")

	return nil
}

//go:embed queries/oauth2_client_tokens/archive_by_refresh.sql
var archiveOAuth2ClientTokenByRefreshQuery string

// ArchiveOAuth2ClientTokenByRefresh archives an OAuth2 client token from the database by its ID.
func (q *Querier) ArchiveOAuth2ClientTokenByRefresh(ctx context.Context, refresh string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if refresh == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientTokenRefreshKey, refresh)
	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenRefreshKey, refresh)

	args := []any{
		refresh,
	}

	if err := q.performWriteQuery(ctx, q.db, "oauth2 client token archive by refresh", archiveOAuth2ClientTokenByRefreshQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating oauth2 client token by refresh")
	}

	logger.Info("oauth2 client token archived by refresh")

	return nil
}
