package tracing

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/require"
)

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

func TestAttachToSpan(t *testing.T) {
	type testCase struct {
		x             any
		name          string
		attachmentKey string
	}
	tests := []testCase{

		{
			name:          "int",
			attachmentKey: t.Name(),
			x:             1,
		},
		{
			name:          "int xlice",
			attachmentKey: t.Name(),
			x:             []int{1},
		},
		{
			name:          "int8",
			attachmentKey: t.Name(),
			x:             int8(1),
		},
		{
			name:          "int8 xlice",
			attachmentKey: t.Name(),
			x:             []int8{int8(1)},
		},
		{
			name:          "int16",
			attachmentKey: t.Name(),
			x:             int16(1),
		},
		{
			name:          "int16 xlice",
			attachmentKey: t.Name(),
			x:             []int16{int16(1)},
		},
		{
			name:          "int32",
			attachmentKey: t.Name(),
			x:             int32(1),
		},
		{
			name:          "int32 xlice",
			attachmentKey: t.Name(),
			x:             []int32{int32(1)},
		},
		{
			name:          "int64",
			attachmentKey: t.Name(),
			x:             int64(1),
		},
		{
			name:          "int64 xlice",
			attachmentKey: t.Name(),
			x:             []int64{int64(1)},
		},
		{
			name:          "uint",
			attachmentKey: t.Name(),
			x:             uint(1),
		},
		{
			name:          "uint xlice",
			attachmentKey: t.Name(),
			x:             []uint{uint(1)},
		},
		{
			name:          "uint8",
			attachmentKey: t.Name(),
			x:             uint8(1),
		},
		{
			name:          "uint8 xlice",
			attachmentKey: t.Name(),
			x:             []uint8{uint8(1)},
		},
		{
			name:          "uint16",
			attachmentKey: t.Name(),
			x:             uint16(1),
		},
		{
			name:          "uint16 xlice",
			attachmentKey: t.Name(),
			x:             []uint16{uint16(1)},
		},
		{
			name:          "uint32",
			attachmentKey: t.Name(),
			x:             uint32(1),
		},
		{
			name:          "uint32 xlice",
			attachmentKey: t.Name(),
			x:             []uint32{uint32(1)},
		},
		{
			name:          "uint64",
			attachmentKey: t.Name(),
			x:             uint64(1),
		},
		{
			name:          "uint64 xlice",
			attachmentKey: t.Name(),
			x:             []uint64{uint64(1)},
		},
		{
			name:          "float32",
			attachmentKey: t.Name(),
			x:             float32(1),
		},
		{
			name:          "float32 xlice",
			attachmentKey: t.Name(),
			x:             []float32{float32(1)},
		},
		{
			name:          "float64",
			attachmentKey: t.Name(),
			x:             float64(1),
		},
		{
			name:          "float64 xlice",
			attachmentKey: t.Name(),
			x:             []float64{float64(1)},
		},
		{
			name:          "string",
			attachmentKey: t.Name(),
			x:             "test",
		},
		{
			name:          "string xlice",
			attachmentKey: t.Name(),
			x:             []string{"test"},
		},
		{
			name:          "bool",
			attachmentKey: t.Name(),
			x:             true,
		},
		{
			name:          "bool xlice",
			attachmentKey: t.Name(),
			x:             []bool{true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, span := StartSpan(context.Background())
			AttachToSpan(span, tt.attachmentKey, tt.x)
		})
	}
}
