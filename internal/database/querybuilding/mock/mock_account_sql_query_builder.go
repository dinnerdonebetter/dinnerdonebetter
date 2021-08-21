package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.HouseholdSQLQueryBuilder = (*HouseholdSQLQueryBuilder)(nil)

// HouseholdSQLQueryBuilder is a mocked types.HouseholdSQLQueryBuilder for testing.
type HouseholdSQLQueryBuilder struct {
	mock.Mock
}

// BuildTransferHouseholdOwnershipQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildTransferHouseholdOwnershipQuery(ctx context.Context, currentOwnerID, newOwnerID, householdID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, currentOwnerID, newOwnerID, householdID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetHouseholdQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildGetHouseholdQuery(ctx context.Context, householdID, userID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, householdID, userID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllHouseholdsCountQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildGetAllHouseholdsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfHouseholdsQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildGetBatchOfHouseholdsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetHouseholdsQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildGetHouseholdsQuery(ctx context.Context, userID uint64, forAdmin bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID, forAdmin, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildHouseholdCreationQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildHouseholdCreationQuery(ctx context.Context, input *types.HouseholdCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateHouseholdQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildUpdateHouseholdQuery(ctx context.Context, input *types.Household) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveHouseholdQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildArchiveHouseholdQuery(ctx context.Context, householdID, userID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, householdID, userID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForHouseholdQuery implements our interface.
func (m *HouseholdSQLQueryBuilder) BuildGetAuditLogEntriesForHouseholdQuery(ctx context.Context, householdID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, householdID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
