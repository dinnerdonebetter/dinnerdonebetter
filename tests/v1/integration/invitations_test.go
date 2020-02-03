package integration

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
)

func checkInvitationEquality(t *testing.T, expected, actual *models.Invitation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Code, actual.Code, "expected Code for ID %d to be %v, but it was %v ", expected.ID, expected.Code, actual.Code)
	assert.Equal(t, expected.Consumed, actual.Consumed, "expected Consumed for ID %d to be %v, but it was %v ", expected.ID, expected.Consumed, actual.Consumed)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyInvitation(t *testing.T) *models.Invitation {
	t.Helper()

	x := &models.InvitationCreationInput{
		Code:     fake.Word(),
		Consumed: fake.Bool(),
	}
	y, err := todoClient.CreateInvitation(context.Background(), x)
	require.NoError(t, err)
	return y
}

func TestInvitations(test *testing.T) {
	test.Parallel()

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create invitation
			expected := &models.Invitation{
				Code:     fake.Word(),
				Consumed: fake.Bool(),
			}
			premade, err := todoClient.CreateInvitation(ctx, &models.InvitationCreationInput{
				Code:     expected.Code,
				Consumed: expected.Consumed,
			})
			checkValueAndError(t, premade, err)

			// Assert invitation equality
			checkInvitationEquality(t, expected, premade)

			// Clean up
			err = todoClient.ArchiveInvitation(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetInvitation(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkInvitationEquality(t, expected, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create invitations
			var expected []*models.Invitation
			for i := 0; i < 5; i++ {
				expected = append(expected, buildDummyInvitation(t))
			}

			// Assert invitation list equality
			actual, err := todoClient.GetInvitations(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Invitations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Invitations),
			)

			// Clean up
			for _, x := range actual.Invitations {
				err = todoClient.ArchiveInvitation(ctx, x.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Fetch invitation
			_, err := todoClient.GetInvitation(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create invitation
			expected := &models.Invitation{
				Code:     fake.Word(),
				Consumed: fake.Bool(),
			}
			premade, err := todoClient.CreateInvitation(ctx, &models.InvitationCreationInput{
				Code:     expected.Code,
				Consumed: expected.Consumed,
			})
			checkValueAndError(t, premade, err)

			// Fetch invitation
			actual, err := todoClient.GetInvitation(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert invitation equality
			checkInvitationEquality(t, expected, actual)

			// Clean up
			err = todoClient.ArchiveInvitation(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			err := todoClient.UpdateInvitation(ctx, &models.Invitation{ID: nonexistentID})
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create invitation
			expected := &models.Invitation{
				Code:     fake.Word(),
				Consumed: fake.Bool(),
			}
			premade, err := todoClient.CreateInvitation(tctx, &models.InvitationCreationInput{
				Code:     fake.Word(),
				Consumed: fake.Bool(),
			})
			checkValueAndError(t, premade, err)

			// Change invitation
			premade.Update(expected.ToInput())
			err = todoClient.UpdateInvitation(ctx, premade)
			assert.NoError(t, err)

			// Fetch invitation
			actual, err := todoClient.GetInvitation(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert invitation equality
			checkInvitationEquality(t, expected, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up
			err = todoClient.ArchiveInvitation(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			tctx := context.Background()
			ctx, span := trace.StartSpan(tctx, t.Name())
			defer span.End()

			// Create invitation
			expected := &models.Invitation{
				Code:     fake.Word(),
				Consumed: fake.Bool(),
			}
			premade, err := todoClient.CreateInvitation(ctx, &models.InvitationCreationInput{
				Code:     expected.Code,
				Consumed: expected.Consumed,
			})
			checkValueAndError(t, premade, err)

			// Clean up
			err = todoClient.ArchiveInvitation(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
