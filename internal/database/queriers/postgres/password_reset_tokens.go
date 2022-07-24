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
	_ types.PasswordResetTokenDataManager = (*SQLQuerier)(nil)
)

// scanPasswordResetToken takes a database Scanner (i.e. *sql.Row) and scans the result into a password reset token struct.
func (q *SQLQuerier) scanPasswordResetToken(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.PasswordResetToken, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.PasswordResetToken{}

	targetVars := []interface{}{
		&x.ID,
		&x.Token,
		&x.ExpiresAt,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.RedeemedOn,
		&x.BelongsToUser,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

const getPasswordResetTokenQuery = `SELECT
	password_reset_tokens.id,
	password_reset_tokens.token,
	password_reset_tokens.expires_at,
	password_reset_tokens.created_on,
	password_reset_tokens.last_updated_on,
	password_reset_tokens.redeemed_on,
	password_reset_tokens.belongs_to_user
FROM password_reset_tokens
WHERE password_reset_tokens.redeemed_on IS NULL
AND extract(epoch from NOW()) < password_reset_tokens.expires_at
AND password_reset_tokens.token = $1
`

// GetPasswordResetTokenByToken fetches a password reset token from the database by its token.
func (q *SQLQuerier) GetPasswordResetTokenByToken(ctx context.Context, token string) (*types.PasswordResetToken, error) {
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

	passwordResetToken, _, _, err := q.scanPasswordResetToken(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning passwordResetToken")
	}

	return passwordResetToken, nil
}

const passwordResetTokenCreationQuery = "INSERT INTO password_reset_tokens (id,token,expires_at,belongs_to_user) VALUES ($1,$2,extract(epoch from (NOW() + (30 * interval '1 minutes'))),$3)"

// CreatePasswordResetToken creates a password reset token in the database.
func (q *SQLQuerier) CreatePasswordResetToken(ctx context.Context, input *types.PasswordResetTokenDatabaseCreationInput) (*types.PasswordResetToken, error) {
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
		CreatedOn:     q.currentTime(),
		BelongsToUser: input.BelongsToUser,
	}

	tracing.AttachPasswordResetTokenIDToSpan(span, x.ID)
	logger.Info("password reset token created")

	return x, nil
}

const archivePasswordResetTokenQuery = "UPDATE password_reset_tokens SET redeemed_on = extract(epoch FROM NOW()) WHERE redeemed_on IS NULL AND id = $1"

// RedeemPasswordResetToken redeems a password reset token from the database by its ID.
func (q *SQLQuerier) RedeemPasswordResetToken(ctx context.Context, passwordResetTokenID string) error {
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

	if err := q.performWriteQuery(ctx, q.db, "password reset token archive", archivePasswordResetTokenQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating password reset token")
	}

	logger.Info("password reset token archived")

	return nil
}
