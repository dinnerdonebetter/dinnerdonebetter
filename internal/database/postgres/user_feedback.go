package postgres

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.UserFeedbackDataManager = (*Querier)(nil)

	// userFeedbacksTableColumns are the columns for the user_feedback table.
	userFeedbacksTableColumns = []string{
		"user_feedback.id",
		"user_feedback.prompt",
		"user_feedback.feedback",
		"user_feedback.rating",
		"user_feedback.context",
		"user_feedback.by_user",
		"user_feedback.created_at",
	}
)

// scanUserFeedback takes a database Scanner (i.e. *sql.Row) and scans the result into a user feedback struct.
func (q *Querier) scanUserFeedback(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.UserFeedback, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.UserFeedback{}
	var rawContext string

	targetVars := []any{
		&x.ID,
		&x.Prompt,
		&x.Feedback,
		&x.Rating,
		&rawContext,
		&x.ByUser,
		&x.CreatedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	if err = json.Unmarshal([]byte(rawContext), &x.Context); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "unmarshalling context")
	}

	return x, filteredCount, totalCount, nil
}

// scanUserFeedbacks takes some database rows and turns them into a slice of user feedback.
func (q *Querier) scanUserFeedbacks(ctx context.Context, rows database.ResultIterator, includeCounts bool) (userFeedbacks []*types.UserFeedback, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanUserFeedback(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		userFeedbacks = append(userFeedbacks, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return userFeedbacks, filteredCount, totalCount, nil
}

// GetUserFeedback fetches a list of user feedback from the database that meet a particular filter.
func (q *Querier) GetUserFeedback(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.UserFeedback], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.UserFeedback]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	} else {
		filter = types.DefaultQueryFilter()
	}

	query, args := q.buildListQuery(ctx, "user_feedback", nil, nil, nil, householdOwnershipColumn, userFeedbacksTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "user feedback", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing user feedback list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanUserFeedbacks(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning user feedback")
	}

	return x, nil
}

//go:embed queries/user_feedback/create.sql
var userFeedbackCreationQuery string

// CreateUserFeedback creates a user feedback in the database.
func (q *Querier) CreateUserFeedback(ctx context.Context, input *types.UserFeedbackDatabaseCreationInput) (*types.UserFeedback, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.UserFeedbackIDKey, input.ID)

	rawEncodedContext, err := json.Marshal(input.Context)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "marshalling context")
	}

	args := []any{
		input.ID,
		input.Prompt,
		input.Feedback,
		input.Rating,
		string(rawEncodedContext),
		input.ByUser,
	}

	// create the user feedback.
	if err = q.performWriteQuery(ctx, q.db, "user feedback creation", userFeedbackCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing user feedback creation query")
	}

	x := &types.UserFeedback{
		ID:        input.ID,
		Context:   input.Context,
		Prompt:    input.Prompt,
		Feedback:  input.Feedback,
		ByUser:    input.ByUser,
		Rating:    input.Rating,
		CreatedAt: q.currentTime(),
	}

	tracing.AttachUserFeedbackIDToSpan(span, x.ID)
	logger.Info("user feedback created")

	return x, nil
}
