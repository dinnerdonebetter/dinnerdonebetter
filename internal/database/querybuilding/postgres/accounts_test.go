package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildGetHouseholdQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT households.id, households.external_id, households.name, households.billing_status, households.contact_email, households.contact_phone, households.payment_processor_customer_id, households.subscription_plan_id, households.created_on, households.last_updated_on, households.archived_on, households.belongs_to_user, household_user_memberships.id, household_user_memberships.belongs_to_user, household_user_memberships.belongs_to_household, household_user_memberships.household_roles, household_user_memberships.default_household, household_user_memberships.created_on, household_user_memberships.last_updated_on, household_user_memberships.archived_on FROM households JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id WHERE households.archived_on IS NULL AND households.belongs_to_user = $1 AND households.id = $2"
		expectedArgs := []interface{}{
			exampleHousehold.BelongsToUser,
			exampleHousehold.ID,
		}
		actualQuery, actualArgs := q.BuildGetHouseholdQuery(ctx, exampleHousehold.ID, exampleUser.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllHouseholdsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(households.id) FROM households WHERE households.archived_on IS NULL"
		actualQuery := q.BuildGetAllHouseholdsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfHouseholdsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT households.id, households.external_id, households.name, households.billing_status, households.contact_email, households.contact_phone, households.payment_processor_customer_id, households.subscription_plan_id, households.created_on, households.last_updated_on, households.archived_on, households.belongs_to_user FROM households WHERE households.id > $1 AND households.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfHouseholdsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetHouseholdsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT households.id, households.external_id, households.name, households.billing_status, households.contact_email, households.contact_phone, households.payment_processor_customer_id, households.subscription_plan_id, households.created_on, households.last_updated_on, households.archived_on, households.belongs_to_user, household_user_memberships.id, household_user_memberships.belongs_to_user, household_user_memberships.belongs_to_household, household_user_memberships.household_roles, household_user_memberships.default_household, household_user_memberships.created_on, household_user_memberships.last_updated_on, household_user_memberships.archived_on, (SELECT COUNT(households.id) FROM households WHERE households.archived_on IS NULL AND households.belongs_to_user = $1) as total_count, (SELECT COUNT(households.id) FROM households WHERE households.archived_on IS NULL AND households.belongs_to_user = $2 AND households.created_on > $3 AND households.created_on < $4 AND households.last_updated_on > $5 AND households.last_updated_on < $6) as filtered_count FROM households JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id WHERE households.archived_on IS NULL AND households.belongs_to_user = $7 AND households.created_on > $8 AND households.created_on < $9 AND households.last_updated_on > $10 AND households.last_updated_on < $11 GROUP BY households.id, household_user_memberships.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			exampleUser.ID,
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.BuildGetHouseholdsQuery(ctx, exampleUser.ID, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateHouseholdQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleHousehold.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO households (external_id,name,billing_status,contact_email,contact_phone,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
		expectedArgs := []interface{}{
			exampleHousehold.ExternalID,
			exampleHousehold.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleHousehold.ContactEmail,
			exampleHousehold.ContactPhone,
			exampleHousehold.BelongsToUser,
		}
		actualQuery, actualArgs := q.BuildHouseholdCreationQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateHouseholdQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE households SET name = $1, contact_email = $2, contact_phone = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $4 AND id = $5"
		expectedArgs := []interface{}{
			exampleHousehold.Name,
			exampleHousehold.ContactEmail,
			exampleHousehold.ContactPhone,
			exampleHousehold.BelongsToUser,
			exampleHousehold.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateHouseholdQuery(ctx, exampleHousehold)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveHouseholdQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE households SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleHousehold.ID,
		}
		actualQuery, actualArgs := q.BuildArchiveHouseholdQuery(ctx, exampleHousehold.ID, exampleUser.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForHouseholdQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleHousehold := fakes.BuildFakeHousehold()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'household_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleHousehold.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForHouseholdQuery(ctx, exampleHousehold.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
