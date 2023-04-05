package types

import (
	"context"
	"testing"
	"time"

	"github.com/prixfixeco/backend/internal/authorization"
	"github.com/prixfixeco/backend/internal/observability/logging"

	"github.com/stretchr/testify/assert"
)

func TestChangeActiveHouseholdInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ChangeActiveHouseholdInput{
			HouseholdID: "123",
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestPASETOCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &PASETOCreationInput{
			ClientID:    t.Name(),
			RequestTime: time.Now().Unix(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestSessionContextData_ToBytes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &SessionContextData{}

		assert.NotEmpty(t, x.ToBytes())
	})
}

func TestSessionContextData_AttachToLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &SessionContextData{
			Requester: RequesterInfo{ServicePermissions: authorization.NewServiceRolePermissionChecker(t.Name())},
		}

		assert.NotNil(t, x.AttachToLogger(logging.NewNoopLogger()))
	})
}

func TestSessionContextData_HouseholdRolePermissionsChecker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &SessionContextData{
			ActiveHouseholdID: t.Name(),
			HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
				t.Name(): authorization.NewHouseholdRolePermissionChecker(t.Name()),
			},
		}

		assert.NotNil(t, x.HouseholdRolePermissionsChecker())
	})
}

func TestSessionContextData_ServiceRolePermissionChecker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &SessionContextData{
			ActiveHouseholdID: t.Name(),
			Requester: RequesterInfo{
				ServicePermissions: authorization.NewServiceRolePermissionChecker(t.Name()),
			},
		}

		assert.NotNil(t, x.ServiceRolePermissionChecker())
	})
}
