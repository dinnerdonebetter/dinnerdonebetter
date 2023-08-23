package postgres

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromHouseholds(includeCounts bool, filteredCount uint64, households ...*types.Household) *sqlmock.Rows {
	columns := append(householdsTableColumns, householdsUserMembershipTableColumns...)

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range households {
		for _, y := range x.Members {
			rowValues := []driver.Value{
				&x.ID,
				&x.Name,
				&x.BillingStatus,
				&x.ContactPhone,
				&x.AddressLine1,
				&x.AddressLine2,
				&x.City,
				&x.State,
				&x.ZipCode,
				&x.Country,
				&x.Latitude,
				&x.Longitude,
				&x.PaymentProcessorCustomerID,
				&x.SubscriptionPlanID,
				&x.CreatedAt,
				&x.LastUpdatedAt,
				&x.ArchivedAt,
				&x.BelongsToUser,
				&y.BelongsToUser.ID,
				&y.BelongsToUser.FirstName,
				&y.BelongsToUser.LastName,
				&y.BelongsToUser.Username,
				&y.BelongsToUser.EmailAddress,
				&y.BelongsToUser.EmailAddressVerifiedAt,
				&y.BelongsToUser.AvatarSrc,
				&y.BelongsToUser.RequiresPasswordChange,
				&y.BelongsToUser.PasswordLastChangedAt,
				&y.BelongsToUser.TwoFactorSecretVerifiedAt,
				&y.BelongsToUser.ServiceRole,
				&y.BelongsToUser.AccountStatus,
				&y.BelongsToUser.AccountStatusExplanation,
				&y.BelongsToUser.Birthday,
				&y.BelongsToUser.CreatedAt,
				&y.BelongsToUser.LastUpdatedAt,
				&y.BelongsToUser.ArchivedAt,
				&y.ID,
				&y.BelongsToUser.ID,
				&y.BelongsToHousehold,
				&y.HouseholdRole,
				&y.DefaultHousehold,
				&y.CreatedAt,
				&x.LastUpdatedAt,
				&y.ArchivedAt,
			}

			if includeCounts {
				rowValues = append(rowValues, filteredCount, len(households))
			}

			exampleRows.AddRow(rowValues...)
		}
	}

	return exampleRows
}

func TestQuerier_GetHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, _ := buildTestClient(t)

		actual, err := c.GetHousehold(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, _ := buildTestClient(t)

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetHouseholds(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholds(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHousehold(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateHousehold(ctx, nil))
	})
}

func TestQuerier_ArchiveHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHousehold(ctx, "", exampleUserID))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHouseholdID, ""))
	})
}
