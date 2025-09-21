package auth

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auth/generated"
)

const (
	resourceTypePasswordResetTokens = "password_reset_tokens"
)

var (
	_ auth.PasswordResetTokenDataManager = (*repository)(nil)
)

// TODO: create AuditLogEntries here

// GetPasswordResetTokenByToken fetches a password reset token from the database by its token.
func (r *repository) GetPasswordResetTokenByToken(ctx context.Context, token string) (*auth.PasswordResetToken, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if token == "" {
		return nil, database.ErrEmptyInputProvided
	}
	tracing.AttachToSpan(span, keys.PasswordResetTokenIDKey, token)

	result, err := r.generatedQuerier.GetPasswordResetToken(ctx, r.db, token)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, r.logger, span, "getting password reset token")
	}

	passwordResetToken := &auth.PasswordResetToken{
		CreatedAt:     result.CreatedAt,
		ExpiresAt:     result.ExpiresAt,
		RedeemedAt:    database.TimePointerFromNullTime(result.RedeemedAt),
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ID:            result.ID,
		Token:         result.Token,
		BelongsToUser: result.BelongsToUser,
	}

	return passwordResetToken, nil
}

// CreatePasswordResetToken creates a password reset token in the database.
func (r *repository) CreatePasswordResetToken(ctx context.Context, input *auth.PasswordResetTokenDatabaseCreationInput) (*auth.PasswordResetToken, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.PasswordResetTokenIDKey, input.ID)
	tracing.AttachToSpan(span, keys.PasswordResetTokenIDKey, input.ID)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "beginning transaction")
	}

	// create the password reset token.
	if err = r.generatedQuerier.CreatePasswordResetToken(ctx, tx, &generated.CreatePasswordResetTokenParams{
		ID:            input.ID,
		Token:         input.Token,
		BelongsToUser: input.BelongsToUser,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing password reset token creation query")
	}

	x := &auth.PasswordResetToken{
		ID:            input.ID,
		Token:         input.Token,
		ExpiresAt:     input.ExpiresAt,
		CreatedAt:     r.CurrentTime(),
		BelongsToUser: input.BelongsToUser,
	}

	if _, err = r.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypePasswordResetTokens,
		RelevantID:   input.ID,
		EventType:    audit.AuditLogEventTypeCreated,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("password reset token created")

	return x, nil
}

// RedeemPasswordResetToken redeems a password reset token from the database by its ID.
func (r *repository) RedeemPasswordResetToken(ctx context.Context, passwordResetTokenID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if passwordResetTokenID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.PasswordResetTokenIDKey, passwordResetTokenID)
	tracing.AttachToSpan(span, keys.PasswordResetTokenIDKey, passwordResetTokenID)

	if err := r.generatedQuerier.RedeemPasswordResetToken(ctx, r.db, passwordResetTokenID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving password reset token")
	}

	return nil
}
