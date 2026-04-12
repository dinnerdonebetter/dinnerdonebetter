package indexing

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity/fakes"
	identitymock "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity/mock"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/testutils"

	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/tracing"
	"github.com/primandproper/platform/reflection"
	textsearch "github.com/primandproper/platform/search/text"
	mocksearch "github.com/primandproper/platform/search/text/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleIndexRequest(T *testing.T) {
	T.Parallel()

	T.Run("user index type", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		identityRepo := &identitymock.RepositoryMock{}
		identityRepo.On(reflection.GetMethodName(identityRepo.GetUser), testutils.ContextMatcher, exampleUser.ID).Return(exampleUser, nil)
		identityRepo.On(reflection.GetMethodName(identityRepo.MarkUserAsIndexed), testutils.ContextMatcher, exampleUser.ID).Return(nil)

		mim := &mocksearch.IndexMock[UserSearchSubset]{
			IndexFunc: func(_ context.Context, _ string, _ any) error { return nil },
		}

		cdi := NewCoreDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			identityRepo,
			mim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleUser.ID,
			IndexType: IndexTypeUsers,
			Delete:    false,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t, identityRepo)
	})

	T.Run("deleting user index type", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		identityRepo := &identitymock.RepositoryMock{}
		identityRepo.On(reflection.GetMethodName(identityRepo.GetUser), testutils.ContextMatcher, exampleUser.ID).Return(exampleUser, nil)

		mim := &mocksearch.IndexMock[UserSearchSubset]{
			DeleteFunc: func(_ context.Context, _ string) error { return nil },
		}

		cdi := NewCoreDataIndexer(
			logger,
			tracing.NewNoopTracerProvider(),
			identityRepo,
			mim,
		)

		indexReq := &textsearch.IndexRequest{
			RowID:     exampleUser.ID,
			IndexType: IndexTypeUsers,
			Delete:    true,
		}

		assert.NoError(t, cdi.HandleIndexRequest(ctx, indexReq))

		mock.AssertExpectationsForObjects(t, identityRepo)
	})
}
