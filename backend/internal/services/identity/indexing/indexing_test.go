package indexing

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	identitymock "github.com/dinnerdonebetter/backend/internal/domain/identity/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/platform/search/text"
	mocksearch "github.com/dinnerdonebetter/backend/internal/platform/search/text/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

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
		identityRepo.On("GetUser", testutils.ContextMatcher, exampleUser.ID).Return(exampleUser, nil)
		identityRepo.On("MarkUserAsIndexed", testutils.ContextMatcher, exampleUser.ID).Return(nil)

		uss := ConvertUserToUserSearchSubset(exampleUser)

		mim := &mocksearch.IndexManager[UserSearchSubset]{}
		mim.On("Index", testutils.ContextMatcher, exampleUser.ID, uss).Return(nil)

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

		mock.AssertExpectationsForObjects(t, identityRepo, mim)
	})

	T.Run("deleting user index type", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := t.Context()
		logger := logging.NewNoopLogger()

		identityRepo := &identitymock.RepositoryMock{}
		identityRepo.On("GetUser", testutils.ContextMatcher, exampleUser.ID).Return(exampleUser, nil)

		mim := &mocksearch.IndexManager[UserSearchSubset]{}
		mim.On("Delete", testutils.ContextMatcher, exampleUser.ID).Return(nil)

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

		mock.AssertExpectationsForObjects(t, identityRepo, mim)
	})
}
