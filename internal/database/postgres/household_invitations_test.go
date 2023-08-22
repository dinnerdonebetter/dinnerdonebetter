package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromHouseholdInvitations(includeCounts bool, filteredCount uint64, householdInvitations ...*types.HouseholdInvitation) *sqlmock.Rows {
	columns := []string{
		"household_invitations.id",
		"households.id",
		"households.name",
		"households.billing_status",
		"households.contact_phone",
		"households.address_line_1",
		"households.address_line_2",
		"households.city",
		"households.state",
		"households.zip_code",
		"households.country",
		"households.latitude",
		"households.longitude",
		"households.payment_processor_customer_id",
		"households.subscription_plan_id",
		"households.created_at",
		"households.last_updated_at",
		"households.archived_at",
		"households.belongs_to_user",
		"household_invitations.to_email",
		"household_invitations.to_user",
		"users.id",
		"users.first_name",
		"users.last_name",
		"users.username",
		"users.email_address",
		"users.email_address_verified_at",
		"users.avatar_src",
		"users.hashed_password",
		"users.requires_password_change",
		"users.password_last_changed_at",
		"users.two_factor_secret",
		"users.two_factor_secret_verified_at",
		"users.service_role",
		"users.user_account_status",
		"users.user_account_status_explanation",
		"users.birthday",
		"users.created_at",
		"users.last_updated_at",
		"users.archived_at",
		"household_invitations.to_name",
		"household_invitations.status",
		"household_invitations.note",
		"household_invitations.status_note",
		"household_invitations.token",
		"household_invitations.expires_at",
		"household_invitations.created_at",
		"household_invitations.last_updated_at",
		"household_invitations.archived_at",
	}

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range householdInvitations {
		rowValues := []driver.Value{
			w.ID,
			w.DestinationHousehold.ID,
			w.DestinationHousehold.Name,
			w.DestinationHousehold.BillingStatus,
			w.DestinationHousehold.ContactPhone,
			w.DestinationHousehold.AddressLine1,
			w.DestinationHousehold.AddressLine2,
			w.DestinationHousehold.City,
			w.DestinationHousehold.State,
			w.DestinationHousehold.ZipCode,
			w.DestinationHousehold.Country,
			w.DestinationHousehold.Latitude,
			w.DestinationHousehold.Longitude,
			w.DestinationHousehold.PaymentProcessorCustomerID,
			w.DestinationHousehold.SubscriptionPlanID,
			w.DestinationHousehold.CreatedAt,
			w.DestinationHousehold.LastUpdatedAt,
			w.DestinationHousehold.ArchivedAt,
			w.DestinationHousehold.BelongsToUser,
			w.ToEmail,
			w.ToUser,
			w.FromUser.ID,
			w.FromUser.FirstName,
			w.FromUser.LastName,
			w.FromUser.Username,
			w.FromUser.EmailAddress,
			w.FromUser.EmailAddressVerifiedAt,
			w.FromUser.AvatarSrc,
			w.FromUser.HashedPassword,
			w.FromUser.RequiresPasswordChange,
			w.FromUser.PasswordLastChangedAt,
			w.FromUser.TwoFactorSecret,
			w.FromUser.TwoFactorSecretVerifiedAt,
			w.FromUser.ServiceRole,
			w.FromUser.AccountStatus,
			w.FromUser.AccountStatusExplanation,
			w.FromUser.Birthday,
			w.FromUser.CreatedAt,
			w.FromUser.LastUpdatedAt,
			w.FromUser.ArchivedAt,
			w.ToName,
			w.Status,
			w.Note,
			w.StatusNote,
			w.Token,
			w.ExpiresAt,
			w.CreatedAt,
			w.LastUpdatedAt,
			w.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(householdInvitations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanHouseholdInvitations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanHouseholdInvitations(ctx, mockRows, false)
		assert.Error(t, err)
	})

	T.Run("logs row closing errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, _, err := q.scanHouseholdInvitations(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_HouseholdInvitationExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.HouseholdInvitationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdInvitationByTokenAndID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.DestinationHousehold.Members = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleHouseholdID, exampleHouseholdInvitation.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByTokenAndIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitationByTokenAndID(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByTokenAndID(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByTokenAndID(ctx, exampleHouseholdID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetHouseholdInvitationByHouseholdAndID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.DestinationHousehold.Members = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleHouseholdID, exampleHouseholdInvitation.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByHouseholdAndIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitationByHouseholdAndID(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByHouseholdAndID(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByHouseholdAndID(ctx, exampleHouseholdID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.DestinationHousehold.Members = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleHouseholdID, exampleHouseholdInvitation.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByHouseholdAndIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdInvitationByHouseholdAndID(ctx, exampleHouseholdID, exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdInvitationByEmailAndToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.DestinationHousehold.Members = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleHouseholdInvitation.ToEmail, exampleHouseholdInvitation.Token}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByEmailAndTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		actual, err := c.GetHouseholdInvitationByEmailAndToken(ctx, exampleHouseholdInvitation.ToEmail, exampleHouseholdInvitation.Token)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByEmailAndToken(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInvitationByEmailAndToken(ctx, exampleHouseholdID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHouseholdInvitation(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_setInvitationStatus(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, _ := buildTestClient(t)

		err := c.setInvitationStatus(ctx, c.db, "", exampleHouseholdInvitation.Note, exampleHouseholdInvitation.Status)
		assert.Error(t, err)
	})
}

func TestSQLQuerier_AcceptHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid invitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleToken := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		err := c.AcceptHouseholdInvitation(ctx, "", exampleToken, t.Name())
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		err := c.AcceptHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, "", exampleHouseholdInvitation.Note)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_attachInvitationsToUser(T *testing.T) {
	T.Parallel()

	T.Run("with invalid email address", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, db := buildTestClient(t)

		err := c.attachInvitationsToUser(ctx, c.db, "", exampleUser.ID)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, db := buildTestClient(t)

		err := c.attachInvitationsToUser(ctx, c.db, exampleUser.EmailAddress, "")
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}
