package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkWebhookEquality(t *testing.T, expected, actual *models.Webhook) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ContentType, actual.ContentType)
	assert.Equal(t, expected.URL, actual.URL)
	assert.Equal(t, expected.Method, actual.Method)
	assert.NotZero(t, actual.CreatedOn)
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func TestWebhooks(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fakemodels.BuildFakeWebhook()
			exampleWebhookInput := fakemodels.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := prixfixeClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, premade)

			// Clean up.
			err = prixfixeClient.ArchiveWebhook(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkWebhookEquality(t, exampleWebhook, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhooks.
			var expected []*models.Webhook
			for i := 0; i < 5; i++ {
				exampleWebhook := fakemodels.BuildFakeWebhook()
				exampleWebhookInput := fakemodels.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
				createdWebhook, err := prixfixeClient.CreateWebhook(ctx, exampleWebhookInput)
				checkValueAndError(t, createdWebhook, err)

				expected = append(expected, createdWebhook)
			}

			// Assert webhook list equality.
			actual, err := prixfixeClient.GetWebhooks(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(t, len(expected) <= len(actual.Webhooks))

			// Clean up.
			for _, webhook := range actual.Webhooks {
				err = prixfixeClient.ArchiveWebhook(ctx, webhook.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Fetch webhook.
			_, err := prixfixeClient.GetWebhook(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fakemodels.BuildFakeWebhook()
			exampleWebhookInput := fakemodels.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := prixfixeClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Fetch webhook.
			actual, err := prixfixeClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, actual)

			// Clean up.
			err = prixfixeClient.ArchiveWebhook(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleWebhook := fakemodels.BuildFakeWebhook()
			exampleWebhook.ID = nonexistentID

			err := prixfixeClient.UpdateWebhook(ctx, exampleWebhook)
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fakemodels.BuildFakeWebhook()
			exampleWebhookInput := fakemodels.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := prixfixeClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Change webhook.
			premade.Name = reverse(premade.Name)
			exampleWebhook.Name = premade.Name
			err = prixfixeClient.UpdateWebhook(ctx, premade)
			assert.NoError(t, err)

			// Fetch webhook.
			actual, err := prixfixeClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up.
			err = prixfixeClient.ArchiveWebhook(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fakemodels.BuildFakeWebhook()
			exampleWebhookInput := fakemodels.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := prixfixeClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Clean up.
			err = prixfixeClient.ArchiveWebhook(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
