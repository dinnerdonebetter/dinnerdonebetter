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

var _ types.OAuth2ClientTokenDataManager = (*Querier)(nil)

// scanOAuth2ClientToken takes a database Scanner (i.e. *sql.Row) and scans the result into an OAuth2 client token struct.
func (q *Querier) scanOAuth2ClientToken(ctx context.Context, scan database.Scanner) (x *types.OAuth2ClientToken, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.OAuth2ClientToken{}

	var (
		codeExpiresAt, accessExpiresAt, refreshExpiresAt time.Time
	)

	targetVars := []any{
		&x.ID,
		&x.ClientID,
		&x.BelongsToUser,
		&x.RedirectURI,
		&x.Scope,
		&x.Code,
		&x.CodeChallenge,
		&x.CodeChallengeMethod,
		&x.CodeCreatedAt,
		&codeExpiresAt,
		&x.Access,
		&x.AccessCreatedAt,
		&accessExpiresAt,
		&x.Refresh,
		&x.RefreshCreatedAt,
		&refreshExpiresAt,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, observability.PrepareError(err, span, "")
	}

	x.CodeExpiresIn = codeExpiresAt.Sub(x.CodeCreatedAt)
	x.AccessExpiresIn = accessExpiresAt.Sub(x.AccessCreatedAt)
	x.RefreshExpiresIn = refreshExpiresAt.Sub(x.RefreshCreatedAt)

	decryptedCode, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, x.Code)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token code")
	}
	x.Code = decryptedCode

	decryptedAccess, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, x.Access)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token access")
	}
	x.Access = decryptedAccess

	decryptedRefresh, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, x.Refresh)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token refresh")
	}
	x.Refresh = decryptedRefresh

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

	encryptedCode, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, code)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token code")
	}

	args := []any{
		encryptedCode,
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
func (q *Querier) GetOAuth2ClientTokenByAccess(ctx context.Context, access string) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if access == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientTokenAccessKey, access)
	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenAccessKey, access)

	encryptedAccess, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, access)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token access")
	}

	args := []any{
		encryptedAccess,
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
func (q *Querier) GetOAuth2ClientTokenByRefresh(ctx context.Context, refresh string) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if refresh == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.OAuth2ClientTokenRefreshKey, refresh)
	tracing.AttachStringToSpan(span, keys.OAuth2ClientTokenRefreshKey, refresh)

	encryptedRefresh, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, refresh)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token access")
	}

	args := []any{
		encryptedRefresh,
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

	encryptedCode, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, input.Code)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token code")
	}

	encryptedAccess, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, input.Access)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token access")
	}

	encryptedRefresh, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, input.Refresh)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token refresh")
	}

	args := []any{
		input.ID,
		input.ClientID,
		input.BelongsToUser,
		input.RedirectURI,
		input.Scope,
		encryptedCode,
		input.CodeChallenge,
		input.CodeChallengeMethod,
		now,
		now.Add(input.CodeExpiresIn),
		encryptedAccess,
		now,
		now.Add(input.AccessExpiresIn),
		encryptedRefresh,
		now,
		now.Add(input.RefreshExpiresIn),
	}

	// create the oauth2 client token.
	if err = q.performWriteQuery(ctx, q.db, "oauth2 client token creation", oauth2ClientTokenCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing oauth2 client token creation query")
	}

	logger.Info("oauth2 client token created")

	oauth2ClientToken := &types.OAuth2ClientToken{
		RefreshCreatedAt:    input.RefreshCreatedAt,
		AccessCreatedAt:     input.AccessCreatedAt,
		CodeCreatedAt:       input.CodeCreatedAt,
		RedirectURI:         input.RedirectURI,
		Scope:               input.Scope,
		Code:                input.Code,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
		BelongsToUser:       input.BelongsToUser,
		Access:              input.Access,
		ClientID:            input.ClientID,
		Refresh:             input.Refresh,
		ID:                  input.ID,
		CodeExpiresIn:       input.CodeExpiresIn,
		AccessExpiresIn:     input.AccessExpiresIn,
		RefreshExpiresIn:    input.RefreshExpiresIn,
	}

	return oauth2ClientToken, nil
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

	encryptedAccess, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, access)
	if err != nil {
		return observability.PrepareError(err, span, "decrypting oauth2 token access")
	}

	args := []any{
		encryptedAccess,
	}

	if err = q.performWriteQuery(ctx, q.db, "oauth2 client token archive by access", archiveOAuth2ClientTokenByAccessQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving oauth2 client token by access")
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

	encryptedCode, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, code)
	if err != nil {
		return observability.PrepareError(err, span, "decrypting oauth2 token access")
	}

	args := []any{
		encryptedCode,
	}

	if err = q.performWriteQuery(ctx, q.db, "oauth2 client token archive by code", archiveOAuth2ClientTokenByCodeQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving oauth2 client token by code")
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

	encryptedRefresh, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, refresh)
	if err != nil {
		return observability.PrepareError(err, span, "decrypting oauth2 token access")
	}

	args := []any{
		encryptedRefresh,
	}

	if err = q.performWriteQuery(ctx, q.db, "oauth2 client token archive by refresh", archiveOAuth2ClientTokenByRefreshQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving oauth2 client token by refresh")
	}

	logger.Info("oauth2 client token archived by refresh")

	return nil
}
