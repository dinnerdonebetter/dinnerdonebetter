package postgres

import (
	"context"
	"testing"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostgres_BuildInvitationExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleInvitation := fakes.BuildFakeInvitation()

		expectedQuery := "SELECT EXISTS ( SELECT invitations.id FROM invitations WHERE invitations.archived_on IS NULL AND invitations.id = $1 )"
		expectedArgs := []interface{}{
			exampleInvitation.ID,
		}
		actualQuery, actualArgs := q.BuildInvitationExistsQuery(ctx, exampleInvitation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleInvitation := fakes.BuildFakeInvitation()

		expectedQuery := "SELECT invitations.id, invitations.external_id, invitations.code, invitations.consumed, invitations.created_on, invitations.last_updated_on, invitations.archived_on, invitations.belongs_to_household FROM invitations WHERE invitations.archived_on IS NULL AND invitations.id = $1"
		expectedArgs := []interface{}{
			exampleInvitation.ID,
		}
		actualQuery, actualArgs := q.BuildGetInvitationQuery(ctx, exampleInvitation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAllInvitationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(invitations.id) FROM invitations WHERE invitations.archived_on IS NULL"
		actualQuery := q.BuildGetAllInvitationsCountQuery(ctx)

		assertArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_BuildGetBatchOfInvitationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT invitations.id, invitations.external_id, invitations.code, invitations.consumed, invitations.created_on, invitations.last_updated_on, invitations.archived_on, invitations.belongs_to_household FROM invitations WHERE invitations.id > $1 AND invitations.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := q.BuildGetBatchOfInvitationsQuery(ctx, beginID, endID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetInvitationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT invitations.id, invitations.external_id, invitations.code, invitations.consumed, invitations.created_on, invitations.last_updated_on, invitations.archived_on, invitations.belongs_to_household, (SELECT COUNT(invitations.id) FROM invitations WHERE invitations.archived_on IS NULL) as total_count, (SELECT COUNT(invitations.id) FROM invitations WHERE invitations.archived_on IS NULL AND invitations.created_on > $1 AND invitations.created_on < $2 AND invitations.last_updated_on > $3 AND invitations.last_updated_on < $4) as filtered_count FROM invitations WHERE invitations.archived_on IS NULL AND invitations.created_on > $5 AND invitations.created_on < $6 AND invitations.last_updated_on > $7 AND invitations.last_updated_on < $8 GROUP BY invitations.id ORDER BY invitations.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.BuildGetInvitationsQuery(ctx, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetInvitationsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT invitations.id, invitations.external_id, invitations.code, invitations.consumed, invitations.created_on, invitations.last_updated_on, invitations.archived_on, invitations.belongs_to_household FROM (SELECT invitations.id, invitations.external_id, invitations.code, invitations.consumed, invitations.created_on, invitations.last_updated_on, invitations.archived_on, invitations.belongs_to_household FROM invitations JOIN unnest('{789,123,456}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT 20) AS invitations WHERE invitations.archived_on IS NULL AND invitations.belongs_to_household = $1 AND invitations.id IN ($2,$3,$4)"
		expectedArgs := []interface{}{
			exampleHouseholdID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := q.BuildGetInvitationsWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs, true)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildCreateInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		exIDGen := &querybuilding.MockExternalIDGenerator{}
		exIDGen.On("NewExternalID").Return(exampleInvitation.ExternalID)
		q.externalIDGenerator = exIDGen

		expectedQuery := "INSERT INTO invitations (external_id,code,consumed,belongs_to_household) VALUES ($1,$2,$3,$4) RETURNING id"
		expectedArgs := []interface{}{
			exampleInvitation.ExternalID,
			exampleInvitation.Code,
			exampleInvitation.Consumed,
			exampleInvitation.BelongsToHousehold,
		}
		actualQuery, actualArgs := q.BuildCreateInvitationQuery(ctx, exampleInput)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)

		mock.AssertExpectationsForObjects(t, exIDGen)
	})
}

func TestPostgres_BuildUpdateInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleInvitation := fakes.BuildFakeInvitation()

		expectedQuery := "UPDATE invitations SET code = $1, consumed = $2, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_household = $3 AND id = $4"
		expectedArgs := []interface{}{
			exampleInvitation.Code,
			exampleInvitation.Consumed,
			exampleInvitation.BelongsToHousehold,
			exampleInvitation.ID,
		}
		actualQuery, actualArgs := q.BuildUpdateInvitationQuery(ctx, exampleInvitation)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildArchiveInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleInvitationID := fakes.BuildFakeID()

		expectedQuery := "UPDATE invitations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"
		expectedArgs := []interface{}{
			exampleInvitationID,
		}
		actualQuery, actualArgs := q.BuildArchiveInvitationQuery(ctx, exampleInvitationID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_BuildGetAuditLogEntriesForInvitationQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestService(t)
		ctx := context.Background()

		exampleInvitation := fakes.BuildFakeInvitation()

		expectedQuery := "SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'invitation_id' = $1 ORDER BY audit_log.created_on"
		expectedArgs := []interface{}{
			exampleInvitation.ID,
		}
		actualQuery, actualArgs := q.BuildGetAuditLogEntriesForInvitationQuery(ctx, exampleInvitation.ID)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
