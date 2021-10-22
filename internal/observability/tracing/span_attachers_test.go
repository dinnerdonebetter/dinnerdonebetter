package tracing

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	useragent "github.com/mssola/user_agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
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

		AttachFilterToSpan(span, 1, 2, t.Name())
	})
}

func TestAttachHouseholdIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachHouseholdIDToSpan(span, "123")
	})
}

func TestAttachRequestingUserIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachRequestingUserIDToSpan(span, "123")
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

func TestAttachAPIClientDatabaseIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachAPIClientDatabaseIDToSpan(span, "123")
	})
}

func TestAttachAPIClientClientIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachAPIClientClientIDToSpan(span, "123")
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

func TestAttachUserIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachUserIDToSpan(span, "123")
	})
}

func TestAttachUsernameToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachUsernameToSpan(span, t.Name())
	})
}

func TestAttachWebhookIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachWebhookIDToSpan(span, "123")
	})
}

func TestAttachURLToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())
		u, err := url.ParseRequestURI("https://todo.verygoodsoftwarenotvirus.ru")
		assert.NoError(t, err)

		AttachURLToSpan(span, u)
	})
}

func TestAttachRequestURIToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachRequestURIToSpan(span, t.Name())
	})
}

func TestAttachRequestToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx, span := StartSpan(context.Background())
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
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

func TestAttachSearchQueryToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachSearchQueryToSpan(span, t.Name())
	})
}

func TestAttachUserAgentDataToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachUserAgentDataToSpan(span, &useragent.UserAgent{})
	})
}

func TestAttachValidInstrumentIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachValidInstrumentIDToSpan(span, "123")
	})
}

func TestAttachValidIngredientIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachValidIngredientIDToSpan(span, "123")
	})
}

func TestAttachValidPreparationIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachValidPreparationIDToSpan(span, "123")
	})
}

func TestAttachValidIngredientPreparationIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachValidIngredientPreparationIDToSpan(span, "123")
	})
}

func TestAttachRecipeIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachRecipeIDToSpan(span, "123")
	})
}

func TestAttachRecipeStepIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachRecipeStepIDToSpan(span, "123")
	})
}

func TestAttachRecipeStepInstrumentIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachRecipeStepInstrumentIDToSpan(span, "123")
	})
}

func TestAttachRecipeStepIngredientIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachRecipeStepIngredientIDToSpan(span, "123")
	})
}

func TestAttachRecipeStepProductIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachRecipeStepProductIDToSpan(span, "123")
	})
}

func TestAttachMealPlanIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachMealPlanIDToSpan(span, "123")
	})
}

func TestAttachMealPlanOptionIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachMealPlanOptionIDToSpan(span, "123")
	})
}

func TestAttachMealPlanOptionVoteIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, span := StartSpan(context.Background())

		AttachMealPlanOptionVoteIDToSpan(span, "123")
	})
}
