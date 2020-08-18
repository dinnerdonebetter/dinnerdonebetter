package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkInvitationEquality(t *testing.T, expected, actual *models.Invitation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Code, actual.Code, "expected Code for ID %d to be %v, but it was %v ", expected.ID, expected.Code, actual.Code)
	assert.Equal(t, expected.Consumed, actual.Consumed, "expected Consumed for ID %d to be %v, but it was %v ", expected.ID, expected.Consumed, actual.Consumed)
	assert.NotZero(t, actual.CreatedOn)
}

func TestInvitations(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create invitation.
			exampleInvitation := fakemodels.BuildFakeInvitation()
			exampleInvitationInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := prixfixeClient.CreateInvitation(ctx, exampleInvitationInput)
			checkValueAndError(t, createdInvitation, err)

			// Assert invitation equality.
			checkInvitationEquality(t, exampleInvitation, createdInvitation)

			// Clean up.
			err = prixfixeClient.ArchiveInvitation(ctx, createdInvitation.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetInvitation(ctx, createdInvitation.ID)
			checkValueAndError(t, actual, err)
			checkInvitationEquality(t, exampleInvitation, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create invitations.
			var expected []*models.Invitation
			for i := 0; i < 5; i++ {
				// Create invitation.
				exampleInvitation := fakemodels.BuildFakeInvitation()
				exampleInvitationInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
				createdInvitation, invitationCreationErr := prixfixeClient.CreateInvitation(ctx, exampleInvitationInput)
				checkValueAndError(t, createdInvitation, invitationCreationErr)

				expected = append(expected, createdInvitation)
			}

			// Assert invitation list equality.
			actual, err := prixfixeClient.GetInvitations(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Invitations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Invitations),
			)

			// Clean up.
			for _, createdInvitation := range actual.Invitations {
				err = prixfixeClient.ArchiveInvitation(ctx, createdInvitation.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent invitation.
			actual, err := prixfixeClient.InvitationExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant invitation exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create invitation.
			exampleInvitation := fakemodels.BuildFakeInvitation()
			exampleInvitationInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := prixfixeClient.CreateInvitation(ctx, exampleInvitationInput)
			checkValueAndError(t, createdInvitation, err)

			// Fetch invitation.
			actual, err := prixfixeClient.InvitationExists(ctx, createdInvitation.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up invitation.
			assert.NoError(t, prixfixeClient.ArchiveInvitation(ctx, createdInvitation.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent invitation.
			_, err := prixfixeClient.GetInvitation(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create invitation.
			exampleInvitation := fakemodels.BuildFakeInvitation()
			exampleInvitationInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := prixfixeClient.CreateInvitation(ctx, exampleInvitationInput)
			checkValueAndError(t, createdInvitation, err)

			// Fetch invitation.
			actual, err := prixfixeClient.GetInvitation(ctx, createdInvitation.ID)
			checkValueAndError(t, actual, err)

			// Assert invitation equality.
			checkInvitationEquality(t, exampleInvitation, actual)

			// Clean up invitation.
			assert.NoError(t, prixfixeClient.ArchiveInvitation(ctx, createdInvitation.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleInvitation := fakemodels.BuildFakeInvitation()
			exampleInvitation.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateInvitation(ctx, exampleInvitation))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create invitation.
			exampleInvitation := fakemodels.BuildFakeInvitation()
			exampleInvitationInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := prixfixeClient.CreateInvitation(ctx, exampleInvitationInput)
			checkValueAndError(t, createdInvitation, err)

			// Change invitation.
			createdInvitation.Update(exampleInvitation.ToUpdateInput())
			err = prixfixeClient.UpdateInvitation(ctx, createdInvitation)
			assert.NoError(t, err)

			// Fetch invitation.
			actual, err := prixfixeClient.GetInvitation(ctx, createdInvitation.ID)
			checkValueAndError(t, actual, err)

			// Assert invitation equality.
			checkInvitationEquality(t, exampleInvitation, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up invitation.
			assert.NoError(t, prixfixeClient.ArchiveInvitation(ctx, createdInvitation.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, prixfixeClient.ArchiveInvitation(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create invitation.
			exampleInvitation := fakemodels.BuildFakeInvitation()
			exampleInvitationInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
			createdInvitation, err := prixfixeClient.CreateInvitation(ctx, exampleInvitationInput)
			checkValueAndError(t, createdInvitation, err)

			// Clean up invitation.
			assert.NoError(t, prixfixeClient.ArchiveInvitation(ctx, createdInvitation.ID))
		})
	})
}
