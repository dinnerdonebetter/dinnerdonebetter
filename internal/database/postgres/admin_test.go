package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestQuerier_UpdateUserAccountStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := &types.UserAccountStatusUpdateInput{
			TargetUserID: exampleUser.ID,
			NewStatus:    "new",
			Reason:       "because",
		}

		ctx := context.Background()
		c, _ := buildTestClient(t)

		args := &generated.SetUserAccountStatusParams{
			UserAccountStatus:            string(exampleInput.NewStatus),
			UserAccountStatusExplanation: exampleInput.Reason,
			ID:                           exampleInput.TargetUserID,
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"SetUserAccountStatus",
			testutils.ContextMatcher,
			args,
		).Return(nil)
		c.generatedQuerier = mockGeneratedQuerier

		assert.NoError(t, c.UpdateUserAccountStatus(ctx, exampleUser.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := &types.UserAccountStatusUpdateInput{
			TargetUserID: exampleUser.ID,
			NewStatus:    "new",
			Reason:       "because",
		}

		ctx := context.Background()
		c, _ := buildTestClient(t)

		args := &generated.SetUserAccountStatusParams{
			UserAccountStatus:            string(exampleInput.NewStatus),
			UserAccountStatusExplanation: exampleInput.Reason,
			ID:                           exampleInput.TargetUserID,
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"SetUserAccountStatus",
			testutils.ContextMatcher,
			args,
		).Return(errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		assert.Error(t, c.UpdateUserAccountStatus(ctx, exampleUser.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})
}
