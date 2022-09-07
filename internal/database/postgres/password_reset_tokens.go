package postgres

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.PasswordResetTokenDataManager = (*Querier)(nil)

	// passwordResetTokensTableColumns are the columns for the password_reset_tokens table.
	passwordResetTokensTableColumns = []string{
		"password_reset_tokens.id",
		"password_reset_tokens.token",
		"password_reset_tokens.expires_at",
		"password_reset_tokens.created_at",
		"password_reset_tokens.last_updated_at",
		"password_reset_tokens.redeemed_at",
		"password_reset_tokens.belongs_to_user",
	}
)

// scanPasswordResetToken takes a database Scanner (i.e. *sql.Row) and scans the result into a password reset token struct.
func (q *Querier) scanPasswordResetToken(ctx context.Context, scan database.Scanner) (x *types.PasswordResetToken, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.PasswordResetToken{}

	targetVars := []interface{}{
		&x.ID,
		&x.Token,
		&x.ExpiresAt,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.RedeemedAt,
		&x.BelongsToUser,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, q.logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

const getPasswordResetTokenQuery = `SELECT
	password_reset_tokens.id,
	password_reset_tokens.token,
	password_reset_tokens.expires_at,
	password_reset_tokens.created_at,
	password_reset_tokens.last_updated_at,
	password_reset_tokens.redeemed_at,
	password_reset_tokens.belongs_to_user
FROM password_reset_tokens
WHERE password_reset_tokens.redeemed_at IS NULL
AND NOW() < password_reset_tokens.expires_at
AND password_reset_tokens.token = $1
`

// GetPasswordResetTokenByToken fetches a password reset token from the database by its token.
func (q *Querier) GetPasswordResetTokenByToken(ctx context.Context, token string) (*types.PasswordResetToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if token == "" {
		return nil, ErrEmptyInputProvided
	}
	tracing.AttachPasswordResetTokenToSpan(span, token)

	args := []interface{}{
		token,
	}

	row := q.getOneRow(ctx, q.db, "passwordResetToken", getPasswordResetTokenQuery, args)

	passwordResetToken, _, _, err := q.scanPasswordResetToken(ctx, row)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning password reset token")
	}

	return passwordResetToken, nil
}

const passwordResetTokenCreationQuery = "INSERT INTO password_reset_tokens (id,token,expires_at,belongs_to_user) VALUES ($1,$2,NOW() + (30 * interval '1 minutes'),$3)"

// CreatePasswordResetToken creates a password reset token in the database.
func (q *Querier) CreatePasswordResetToken(ctx context.Context, input *types.PasswordResetTokenDatabaseCreationInput) (*types.PasswordResetToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.PasswordResetTokenIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Token,
		input.BelongsToUser,
	}

	// create the password reset token.
	if err := q.performWriteQuery(ctx, q.db, "password reset token creation", passwordResetTokenCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing password reset token creation query")
	}

	x := &types.PasswordResetToken{
		ID:            input.ID,
		Token:         input.Token,
		ExpiresAt:     input.ExpiresAt,
		CreatedAt:     q.currentTime(),
		BelongsToUser: input.BelongsToUser,
	}

	tracing.AttachPasswordResetTokenIDToSpan(span, x.ID)
	logger.Info("password reset token created")

	return x, nil
}

const redeemPasswordResetTokenQuery = "UPDATE password_reset_tokens SET redeemed_at = NOW() WHERE redeemed_at IS NULL AND id = $1"

// RedeemPasswordResetToken redeems a password reset token from the database by its ID.
func (q *Querier) RedeemPasswordResetToken(ctx context.Context, passwordResetTokenID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if passwordResetTokenID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.PasswordResetTokenIDKey, passwordResetTokenID)
	tracing.AttachPasswordResetTokenIDToSpan(span, passwordResetTokenID)

	args := []interface{}{
		passwordResetTokenID,
	}

	if err := q.performWriteQuery(ctx, q.db, "password reset token archive", redeemPasswordResetTokenQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating password reset token")
	}

	logger.Info("password reset token archived")

	return nil
}
