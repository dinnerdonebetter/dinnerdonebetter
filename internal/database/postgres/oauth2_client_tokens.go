package postgres

import (
	"context"
	_ "embed"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// var _ types.OAuth2ClientTokenDataManager = (*Querier)(nil)

// scanOAuth2ClientToken takes a database Scanner (i.e. *sql.Row) and scans the result into an OAuth2 client token struct.
func (q *Querier) scanOAuth2ClientToken(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.OAuth2ClientToken, filteredCount, totalCount uint64, err error) {
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

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, filteredCount, totalCount, nil
}

//go:embed queries/oauth2_client_tokens/get_one_by_access.sql
var getOAuth2ClientTokenQuery string

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

	row := q.getOneRow(ctx, q.db, "oauth2ClientToken", getOAuth2ClientTokenQuery, args)

	oauth2ClientToken, _, _, err := q.scanOAuth2ClientToken(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning oauth2ClientToken")
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

	args := []any{
		input.ID,
		input.ClientID,
		input.BelongsToUser,
		input.RedirectURI,
		input.Scope,
		input.Code,
		input.CodeChallenge,
		input.CodeChallengeMethod,
		input.CodeCreateAt,
		time.Now().Add(input.CodeExpiresIn),
		input.Access,
		input.AccessCreateAt,
		time.Now().Add(input.AccessExpiresIn),
		input.Refresh,
		input.RefreshCreateAt,
		time.Now().Add(input.RefreshExpiresIn),
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
