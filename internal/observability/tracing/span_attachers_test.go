package tracing

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func Test_attachUint8ToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		attachUint8ToSpan(span, t.Name(), 1)
	})
}

func Test_attachUint64ToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		attachUint64ToSpan(span, t.Name(), 123)
	})
}

func Test_attachStringToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		attachStringToSpan(span, t.Name(), "blah")
	})
}

func Test_attachBooleanToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		attachBooleanToSpan(span, t.Name(), false)
	})
}

func Test_attachSliceToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		attachSliceToSpan(span, t.Name(), []string{t.Name()})
	})
}

func TestAttachToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachToSpan(span, t.Name(), "blah")
	})
}

func TestAttachFilterToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachFilterDataToSpan(span, func(x uint64) *uint64 { return &x }(1), func(x uint8) *uint8 { return &x }(2), func(x string) *string { return &x }(t.Name()))
	})
}

func TestAttachSessionContextDataToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachSessionContextDataToSpan(span, &types.SessionContextData{
			HouseholdPermissions: nil,
			Requester: types.RequesterInfo{
				ServicePermissions: authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
			},
			ActiveHouseholdID: "",
		})
	})
}

func TestAttachUserToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		_, span := StartSpan(context.Background())

		AttachUserToSpan(span, exampleUser)
	})
}

func TestAttachRequestToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx, span := StartSpan(context.Background())
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", http.NoBody)
		req.Header.Set(t.Name(), "blah")
		require.NoError(t, err)

		AttachRequestToSpan(span, req)
	})
}

func TestAttachResponseToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())
		res := &http.Response{
			Header: map[string][]string{},
		}
		res.Header.Set(t.Name(), "blah")

		AttachResponseToSpan(span, res)
	})
}

func TestAttachErrorToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachErrorToSpan(span, t.Name(), errors.New("blah"))
	})
}

func TestAttachDatabaseQueryToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachDatabaseQueryToSpan(span, "description", "query", []interface{}{"blah"})
	})
}

func TestAttachQueryFilterToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachQueryFilterToSpan(span, types.DefaultQueryFilter())
	})

	T.Run("with nil", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachQueryFilterToSpan(span, nil)
	})
}
