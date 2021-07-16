package mock

import (
	"context"

	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.InvitationSQLQueryBuilder = (*InvitationSQLQueryBuilder)(nil)

// InvitationSQLQueryBuilder is a mocked types.InvitationSQLQueryBuilder for testing.
type InvitationSQLQueryBuilder struct {
	mock.Mock
}

// BuildInvitationExistsQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildInvitationExistsQuery(ctx context.Context, invitationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, invitationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetInvitationQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildGetInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, invitationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAllInvitationsCountQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildGetAllInvitationsCountQuery(ctx context.Context) string {
	returnArgs := m.Called(ctx)

	return returnArgs.String(0)
}

// BuildGetBatchOfInvitationsQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildGetBatchOfInvitationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, beginID, endID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetInvitationsQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildGetInvitationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, includeArchived, filter)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetInvitationsWithIDsQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildGetInvitationsWithIDsQuery(ctx context.Context, accountID uint64, limit uint8, ids []uint64, restrictToAccount bool) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, accountID, limit, ids, restrictToAccount)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateInvitationQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildCreateInvitationQuery(ctx context.Context, input *types.InvitationCreationInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetAuditLogEntriesForInvitationQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildGetAuditLogEntriesForInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, invitationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUpdateInvitationQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildUpdateInvitationQuery(ctx context.Context, input *types.Invitation) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveInvitationQuery implements our interface.
func (m *InvitationSQLQueryBuilder) BuildArchiveInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, invitationID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
