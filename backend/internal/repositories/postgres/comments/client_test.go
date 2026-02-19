package comments

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/comments"
	"github.com/dinnerdonebetter/backend/internal/domain/comments/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/migrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pgcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func buildDatabaseClientForTest(t *testing.T) (*repository, *pgcontainers.PostgresContainer) {
	t.Helper()

	ctx := t.Context()
	container, db, config := pgtesting.BuildDatabaseContainerForTest(t)
	require.NoError(t, migrations.NewMigrator(logging.NewNoopLogger()).Migrate(ctx, db))

	pgc, err := postgres.ProvideDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), config)
	require.NotNil(t, pgc)
	require.NoError(t, err)

	auditLogEntryRepo := auditlogentries.ProvideAuditLogRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), pgc)

	c := ProvideCommentsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, pgc)

	return c.(*repository), container
}

func buildInertClientForTest(t *testing.T) *repository {
	t.Helper()

	c := ProvideCommentsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), nil, &database.MockClient{})

	return c.(*repository)
}

func createCommentForTest(t *testing.T, ctx context.Context, input *comments.CommentDatabaseCreationInput, dbc *repository) *comments.Comment {
	t.Helper()

	if input == nil {
		input = fakes.BuildFakeCommentDatabaseCreationInput()
	}

	created, err := dbc.CreateComment(ctx, input)
	assert.NoError(t, err)
	require.NotNil(t, created)

	fetched, err := dbc.GetComment(ctx, created.ID)
	assert.NoError(t, err)
	require.NotNil(t, fetched)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.Content, fetched.Content)
	assert.Equal(t, created.TargetType, fetched.TargetType)
	assert.Equal(t, created.ReferencedID, fetched.ReferencedID)
	assert.Equal(t, created.BelongsToUser, fetched.BelongsToUser)

	return created
}

func TestQuerier_Integration_Comments(t *testing.T) {
	if !pgtesting.RunContainerTests {
		t.SkipNow()
	}

	ctx := t.Context()
	dbc, container := buildDatabaseClientForTest(t)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := pgtesting.CreateUserForTest(t, nil, dbc.writeDB)
	referencedID := identifiers.New()
	targetType := comments.CommentTargetTypeRecipes

	input := fakes.BuildFakeCommentDatabaseCreationInput()
	input.BelongsToUser = user.ID
	input.ReferencedID = referencedID
	input.TargetType = targetType

	// create
	created := createCommentForTest(t, ctx, input, dbc)

	// fetch as list
	result, err := dbc.GetCommentsForReference(ctx, targetType, referencedID, nil)
	assert.NoError(t, err)
	require.NotEmpty(t, result.Data)
	assert.Equal(t, 1, len(result.Data))
	assert.Equal(t, created.ID, result.Data[0].ID)

	// update
	newContent := "updated content"
	err = dbc.UpdateComment(ctx, created.ID, user.ID, newContent)
	assert.NoError(t, err)

	updated, err := dbc.GetComment(ctx, created.ID)
	assert.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, newContent, updated.Content)
	assert.NotNil(t, updated.LastUpdatedAt)

	// archive
	err = dbc.ArchiveComment(ctx, created.ID)
	assert.NoError(t, err)

	fetchedAfterArchive, err := dbc.GetComment(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, fetchedAfterArchive)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestQuerier_Integration_Comments_WithReplies(t *testing.T) {
	if !pgtesting.RunContainerTests {
		t.SkipNow()
	}

	ctx := t.Context()
	dbc, container := buildDatabaseClientForTest(t)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := pgtesting.CreateUserForTest(t, nil, dbc.writeDB)
	referencedID := identifiers.New()
	targetType := comments.CommentTargetTypeRecipes

	// create parent comment
	parentInput := fakes.BuildFakeCommentDatabaseCreationInput()
	parentInput.BelongsToUser = user.ID
	parentInput.ReferencedID = referencedID
	parentInput.TargetType = targetType
	parent := createCommentForTest(t, ctx, parentInput, dbc)

	// create reply
	replyInput := fakes.BuildFakeCommentDatabaseCreationInput()
	replyInput.BelongsToUser = user.ID
	replyInput.ReferencedID = referencedID
	replyInput.TargetType = targetType
	replyInput.ParentCommentID = &parent.ID
	reply := createCommentForTest(t, ctx, replyInput, dbc)

	// fetch all for reference - should get both parent and reply
	result, err := dbc.GetCommentsForReference(ctx, targetType, referencedID, nil)
	assert.NoError(t, err)
	require.Len(t, result.Data, 2)
	var foundParent, foundReply bool
	for _, c := range result.Data {
		if c.ID == parent.ID {
			foundParent = true
			assert.Nil(t, c.ParentCommentID)
		}
		if c.ID == reply.ID {
			foundReply = true
			require.NotNil(t, c.ParentCommentID)
			assert.Equal(t, parent.ID, *c.ParentCommentID)
		}
	}
	assert.True(t, foundParent)
	assert.True(t, foundReply)
}

func TestQuerier_Integration_ArchiveCommentsForReference(t *testing.T) {
	if !pgtesting.RunContainerTests {
		t.SkipNow()
	}

	ctx := t.Context()
	dbc, container := buildDatabaseClientForTest(t)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := pgtesting.CreateUserForTest(t, nil, dbc.writeDB)
	referencedID := identifiers.New()
	targetType := comments.CommentTargetTypeMealPlans

	input := fakes.BuildFakeCommentDatabaseCreationInput()
	input.BelongsToUser = user.ID
	input.ReferencedID = referencedID
	input.TargetType = targetType

	created := createCommentForTest(t, ctx, input, dbc)

	// archive all for reference
	err := dbc.ArchiveCommentsForReference(ctx, targetType, referencedID)
	assert.NoError(t, err)

	// comment should no longer be fetchable (GetComment returns archived)
	fetched, err := dbc.GetComment(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, fetched)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	// list should be empty when not including archived
	result, err := dbc.GetCommentsForReference(ctx, targetType, referencedID, nil)
	assert.NoError(t, err)
	assert.Empty(t, result.Data)
}

func TestQuerier_CreateComment(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.CreateComment(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetComment(T *testing.T) {
	T.Parallel()

	T.Run("with empty id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetComment(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetCommentsForReference(T *testing.T) {
	T.Parallel()

	T.Run("with empty target type", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetCommentsForReference(ctx, "", "ref-id", nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with empty referenced id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetCommentsForReference(ctx, comments.CommentTargetTypeRecipes, "", nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateComment(T *testing.T) {
	T.Parallel()

	T.Run("with empty id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		err := c.UpdateComment(ctx, "", "user-id", "content")
		assert.Error(t, err)
	})

	T.Run("with empty belongs to user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		err := c.UpdateComment(ctx, "comment-id", "", "content")
		assert.Error(t, err)
	})
}

func TestQuerier_ArchiveComment(T *testing.T) {
	T.Parallel()

	T.Run("with empty id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		err := c.ArchiveComment(ctx, "")
		assert.Error(t, err)
	})
}

func TestQuerier_ArchiveCommentsForReference(T *testing.T) {
	T.Parallel()

	T.Run("with empty target type", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		err := c.ArchiveCommentsForReference(ctx, "", "ref-id")
		assert.Error(t, err)
	})

	T.Run("with empty referenced id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		err := c.ArchiveCommentsForReference(ctx, comments.CommentTargetTypeRecipes, "")
		assert.Error(t, err)
	})
}

func TestQuerier_Integration_Comments_CursorPagination(t *testing.T) {
	if !pgtesting.RunContainerTests {
		t.SkipNow()
	}

	ctx := t.Context()
	dbc, container := buildDatabaseClientForTest(t)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := pgtesting.CreateUserForTest(t, nil, dbc.writeDB)
	referencedID := identifiers.New()
	targetType := comments.CommentTargetTypeMeals

	pgtesting.TestCursorBasedPagination(t, ctx, pgtesting.PaginationTestConfig[comments.Comment]{
		TotalItems: 9,
		PageSize:   3,
		ItemName:   "comment",
		CreateItem: func(ctx context.Context, i int) *comments.Comment {
			input := fakes.BuildFakeCommentDatabaseCreationInput()
			input.BelongsToUser = user.ID
			input.ReferencedID = referencedID
			input.TargetType = targetType
			return createCommentForTest(t, ctx, input, dbc)
		},
		FetchPage: func(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[comments.Comment], error) {
			return dbc.GetCommentsForReference(ctx, targetType, referencedID, filter)
		},
		GetID: func(c *comments.Comment) string {
			return c.ID
		},
		CleanupItem: func(ctx context.Context, c *comments.Comment) error {
			return dbc.ArchiveComment(ctx, c.ID)
		},
	})
}
