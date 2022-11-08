package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromHouseholdInvitations(includeCounts bool, filteredCount uint64, householdInvitations ...*types.HouseholdInvitation) *sqlmock.Rows {
	columns := householdInvitationsTableColumns

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
			w.DestinationHousehold.ContactEmail,
			w.DestinationHousehold.ContactPhone,
			w.DestinationHousehold.PaymentProcessorCustomerID,
			w.DestinationHousehold.SubscriptionPlanID,
			w.DestinationHousehold.TimeZone,
			w.DestinationHousehold.CreatedAt,
			w.DestinationHousehold.LastUpdatedAt,
			w.DestinationHousehold.ArchivedAt,
			w.DestinationHousehold.BelongsToUser,
			w.ToEmail,
			w.ToUser,
			w.FromUser.ID,
			w.FromUser.Username,
			w.FromUser.EmailAddress,
			w.FromUser.AvatarSrc,
			w.FromUser.HashedPassword,
			w.FromUser.RequiresPasswordChange,
			w.FromUser.PasswordLastChangedAt,
			w.FromUser.TwoFactorSecret,
			w.FromUser.TwoFactorSecretVerifiedAt,
			w.FromUser.ServiceRole,
			w.FromUser.AccountStatus,
			w.FromUser.AccountStatusExplanation,
			w.FromUser.BirthDay,
			w.FromUser.BirthMonth,
			w.FromUser.CreatedAt,
			w.FromUser.LastUpdatedAt,
			w.FromUser.ArchivedAt,
			w.Status,
			w.Note,
			w.StatusNote,
			w.Token,
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

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(householdInvitationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdInvitationID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.HouseholdInvitationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(householdInvitationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdInvitationID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitationID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdInvitationID,
		}

		db.ExpectQuery(formatQueryForSQLMock(householdInvitationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.HouseholdInvitationExists(ctx, exampleHouseholdInvitationID)
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

		args := []interface{}{exampleHouseholdID, exampleHouseholdInvitation.ID}

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

		args := []interface{}{exampleHouseholdID, exampleHouseholdInvitation.ID}

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
}

func TestQuerier_GetHouseholdInvitationByEmailAndToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.DestinationHousehold.Members = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{exampleHouseholdInvitation.ToEmail, exampleHouseholdInvitation.Token}

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

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.StatusNote = ""
		exampleHouseholdInvitation.DestinationHousehold = types.Household{ID: exampleHouseholdInvitation.DestinationHousehold.ID}
		exampleHouseholdInvitation.FromUser = types.User{ID: exampleHouseholdInvitation.FromUser.ID}

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationDatabaseCreationInput(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.FromUser,
			exampleInput.ToUser,
			exampleInput.Note,
			exampleInput.ToEmail,
			exampleInput.Token,
			exampleInput.DestinationHouseholdID,
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleHouseholdInvitation.CreatedAt
		}

		actual, err := c.CreateHouseholdInvitation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitation, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHouseholdInvitation(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing creation query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.StatusNote = ""
		exampleHouseholdInvitation.DestinationHousehold = types.Household{ID: exampleHouseholdInvitation.DestinationHousehold.ID}
		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationDatabaseCreationInput(exampleHouseholdInvitation)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.FromUser,
			exampleInput.ToUser,
			exampleInput.Note,
			exampleInput.ToEmail,
			exampleInput.Token,
			exampleInput.DestinationHouseholdID,
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdInvitationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleHouseholdInvitation.CreatedAt
		}

		actual, err := c.CreateHouseholdInvitation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_BuildGetPendingHouseholdInvitationsFromUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		expectedQuery := "SELECT household_invitations.id, households.id, households.name, households.billing_status, households.contact_email, households.contact_phone, households.payment_processor_customer_id, households.subscription_plan_id, households.time_zone, households.created_at, households.last_updated_at, households.archived_at, households.belongs_to_user, household_invitations.to_email, household_invitations.to_user, users.id, users.username, users.email_address, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_at, users.two_factor_secret, users.two_factor_secret_verified_at, users.service_roles, users.user_account_status, users.user_account_status_explanation, users.birth_day, users.birth_month, users.created_at, users.last_updated_at, users.archived_at, household_invitations.status, household_invitations.note, household_invitations.status_note, household_invitations.token, household_invitations.created_at, household_invitations.last_updated_at, household_invitations.archived_at, (SELECT COUNT(household_invitations.id) FROM household_invitations JOIN households ON household_invitations.destination_household = households.id JOIN users ON household_invitations.from_user = users.id WHERE household_invitations.archived_at IS NULL AND household_invitations.from_user = $1 AND household_invitations.status = $2) as filtered_count, (SELECT COUNT(household_invitations.id) FROM household_invitations JOIN households ON household_invitations.destination_household = households.id JOIN users ON household_invitations.from_user = users.id WHERE household_invitations.archived_at IS NULL AND household_invitations.from_user = $3 AND household_invitations.status = $4) as total_count FROM household_invitations JOIN households ON household_invitations.destination_household = households.id JOIN users ON household_invitations.from_user = users.id WHERE household_invitations.archived_at IS NULL AND household_invitations.from_user = $5 AND household_invitations.status = $6 LIMIT 20"
		expectedArgs := []interface{}{
			userID,
			types.PendingHouseholdInvitationStatus,
			userID,
			types.PendingHouseholdInvitationStatus,
			userID,
			types.PendingHouseholdInvitationStatus,
		}

		c, _ := buildTestClient(t)

		actualQuery, actualArgs := c.BuildGetPendingHouseholdInvitationsFromUserQuery(ctx, userID, filter)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSQLQuerier_GetPendingHouseholdInvitationsFromUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		for i := range exampleHouseholdInvitationList.HouseholdInvitations {
			exampleHouseholdInvitationList.HouseholdInvitations[i].DestinationHousehold.Members = nil
		}

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsFromUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(true, exampleHouseholdInvitationList.FilteredCount, exampleHouseholdInvitationList.HouseholdInvitations...))

		actual, err := c.GetPendingHouseholdInvitationsFromUser(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		for i := range exampleHouseholdInvitationList.HouseholdInvitations {
			exampleHouseholdInvitationList.HouseholdInvitations[i].DestinationHousehold.Members = nil
		}

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsFromUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetPendingHouseholdInvitationsFromUser(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_BuildGetPendingHouseholdInvitationsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		expectedQuery := "SELECT household_invitations.id, households.id, households.name, households.billing_status, households.contact_email, households.contact_phone, households.payment_processor_customer_id, households.subscription_plan_id, households.time_zone, households.created_at, households.last_updated_at, households.archived_at, households.belongs_to_user, household_invitations.to_email, household_invitations.to_user, users.id, users.username, users.email_address, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_at, users.two_factor_secret, users.two_factor_secret_verified_at, users.service_roles, users.user_account_status, users.user_account_status_explanation, users.birth_day, users.birth_month, users.created_at, users.last_updated_at, users.archived_at, household_invitations.status, household_invitations.note, household_invitations.status_note, household_invitations.token, household_invitations.created_at, household_invitations.last_updated_at, household_invitations.archived_at, (SELECT COUNT(household_invitations.id) FROM household_invitations JOIN households ON household_invitations.destination_household = households.id JOIN users ON household_invitations.from_user = users.id WHERE household_invitations.archived_at IS NULL AND household_invitations.status = $1 AND household_invitations.to_user = $2) as filtered_count, (SELECT COUNT(household_invitations.id) FROM household_invitations JOIN households ON household_invitations.destination_household = households.id JOIN users ON household_invitations.from_user = users.id WHERE household_invitations.archived_at IS NULL AND household_invitations.status = $3 AND household_invitations.to_user = $4) as total_count FROM household_invitations JOIN households ON household_invitations.destination_household = households.id JOIN users ON household_invitations.from_user = users.id WHERE household_invitations.archived_at IS NULL AND household_invitations.status = $5 AND household_invitations.to_user = $6 LIMIT 20"
		expectedArgs := []interface{}{
			types.PendingHouseholdInvitationStatus,
			userID,
			types.PendingHouseholdInvitationStatus,
			userID,
			types.PendingHouseholdInvitationStatus,
			userID,
		}

		c, _ := buildTestClient(t)

		actualQuery, actualArgs := c.BuildGetPendingHouseholdInvitationsForUserQuery(ctx, userID, filter)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSQLQuerier_GetPendingHouseholdInvitationsForUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()
		for i := range exampleHouseholdInvitationList.HouseholdInvitations {
			exampleHouseholdInvitationList.HouseholdInvitations[i].DestinationHousehold.Members = nil
		}

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsForUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(true, exampleHouseholdInvitationList.FilteredCount, exampleHouseholdInvitationList.HouseholdInvitations...))

		actual, err := c.GetPendingHouseholdInvitationsForUser(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInvitationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.BuildGetPendingHouseholdInvitationsForUserQuery(ctx, exampleUserID, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetPendingHouseholdInvitationsForUser(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_setInvitationStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdInvitation.Status,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.setInvitationStatus(ctx, c.db, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note, exampleHouseholdInvitation.Status)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household invitation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, _ := buildTestClient(t)

		err := c.setInvitationStatus(ctx, c.db, "", exampleHouseholdInvitation.Note, exampleHouseholdInvitation.Status)
		assert.Error(t, err)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdInvitation.Status,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		err := c.setInvitationStatus(ctx, c.db, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note, exampleHouseholdInvitation.Status)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_CancelHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			types.CancelledHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.CancelHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_AcceptHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []interface{}{
			types.AcceptedHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		invitationLookupArgs := []interface{}{
			exampleHouseholdInvitation.Token,
			exampleHouseholdInvitation.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByTokenAndIDQuery)).
			WithArgs(interfaceToDriverValue(invitationLookupArgs)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		addUserToHouseholdArgs := []interface{}{
			&idMatcher{},
			*exampleHouseholdInvitation.ToUser,
			exampleHouseholdInvitation.DestinationHousehold.ID,
			"household_member",
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		err := c.AcceptHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, exampleHouseholdInvitation.Note)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		err := c.AcceptHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, exampleHouseholdInvitation.Note)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error setting invitation status", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []interface{}{
			types.AcceptedHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		err := c.AcceptHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, exampleHouseholdInvitation.Note)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error getting household invitation", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []interface{}{
			types.AcceptedHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		invitationLookupArgs := []interface{}{
			exampleHouseholdInvitation.Token,
			exampleHouseholdInvitation.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByTokenAndIDQuery)).
			WithArgs(interfaceToDriverValue(invitationLookupArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		err := c.AcceptHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, exampleHouseholdInvitation.Note)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error adding user to household", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []interface{}{
			types.AcceptedHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		invitationLookupArgs := []interface{}{
			exampleHouseholdInvitation.Token,
			exampleHouseholdInvitation.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByTokenAndIDQuery)).
			WithArgs(interfaceToDriverValue(invitationLookupArgs)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		addUserToHouseholdArgs := []interface{}{
			&idMatcher{},
			*exampleHouseholdInvitation.ToUser,
			exampleHouseholdInvitation.DestinationHousehold.ID,
			"household_member",
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		err := c.AcceptHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, exampleHouseholdInvitation.Note)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []interface{}{
			types.AcceptedHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		invitationLookupArgs := []interface{}{
			exampleHouseholdInvitation.Token,
			exampleHouseholdInvitation.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByTokenAndIDQuery)).
			WithArgs(interfaceToDriverValue(invitationLookupArgs)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		addUserToHouseholdArgs := []interface{}{
			&idMatcher{},
			*exampleHouseholdInvitation.ToUser,
			exampleHouseholdInvitation.DestinationHousehold.ID,
			"household_member",
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		err := c.AcceptHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Token, exampleHouseholdInvitation.Note)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_RejectHouseholdInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		c, db := buildTestClient(t)

		args := []interface{}{
			types.RejectedHouseholdInvitationStatus,
			exampleHouseholdInvitation.Note,
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.RejectHouseholdInvitation(ctx, exampleHouseholdInvitation.ID, exampleHouseholdInvitation.Note)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_attachInvitationsToUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUser.ID,
			exampleUser.EmailAddress,
		}

		db.ExpectExec(formatQueryForSQLMock(attachInvitationsToUserIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.attachInvitationsToUser(ctx, c.db, exampleUser.EmailAddress, exampleUser.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

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
