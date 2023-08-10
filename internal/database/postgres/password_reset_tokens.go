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

	targetVars := []any{
		&x.ID,
		&x.Token,
		&x.ExpiresAt,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.RedeemedAt,
		&x.BelongsToUser,
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	return x, filteredCount, totalCount, nil
}

//go:embed queries/password_reset_tokens/get_one.sql
var getPasswordResetTokenQuery string

// GetPasswordResetTokenByToken fetches a password reset token from the database by its token.
func (q *Querier) GetPasswordResetTokenByToken(ctx context.Context, token string) (*types.PasswordResetToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if token == "" {
		return nil, ErrEmptyInputProvided
	}
	tracing.AttachPasswordResetTokenToSpan(span, token)

	args := []any{
		token,
	}

	row := q.getOneRow(ctx, q.db, "passwordResetToken", getPasswordResetTokenQuery, args)

	passwordResetToken, _, _, err := q.scanPasswordResetToken(ctx, row)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning password reset token")
	}

	return passwordResetToken, nil
}

//go:embed queries/password_reset_tokens/create.sql
var passwordResetTokenCreationQuery string

// CreatePasswordResetToken creates a password reset token in the database.
func (q *Querier) CreatePasswordResetToken(ctx context.Context, input *types.PasswordResetTokenDatabaseCreationInput) (*types.PasswordResetToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.PasswordResetTokenIDKey, input.ID)

	args := []any{
		input.ID,
		input.Token,
		input.BelongsToUser,
	}

	// create the password reset token.
	if err := q.performWriteQuery(ctx, q.db, "password reset token creation", passwordResetTokenCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing password reset token creation query")
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

	if err := q.generatedQuerier.RedeemPasswordResetToken(ctx, q.db, passwordResetTokenID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving password reset token")
	}

	logger.Info("password reset token archived")

	return nil
}
