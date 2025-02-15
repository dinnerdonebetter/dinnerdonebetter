package authentication

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	mockmetrics "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticationService_AuthorizationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
			},
			ActiveHouseholdID:    helper.exampleHousehold.ID,
			HouseholdPermissions: helper.examplePermCheckers,
		}

		mockUserDataManager := &mocktypes.UserDataManagerMock{}
		mockUserDataManager.On(
			"GetSessionContextDataForUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(sessionCtxData, nil)
		helper.service.userDataManager = mockUserDataManager

		h := &testutils.MockHTTPHandler{}
		h.On(
			"ServeHTTP",
			testutils.HTTPResponseWriterMatcher,
			testutils.HTTPRequestMatcher,
		).Return()

		helper.req = helper.req.WithContext(context.WithValue(helper.ctx, sessions.SessionContextDataKey, sessionCtxData))

		helper.service.AuthorizationMiddleware(h).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, h)
	})

	T.Run("with banned user", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.exampleUser.AccountStatus = string(types.BannedUserAccountStatus)
		helper.setContextFetcher(t)

		mp := &mockmetrics.Int64Counter{}
		mp.On("Add", testutils.ContextMatcher, int64(1), mock.Anything).Return()
		helper.service.rejectedRequestCounter = mp

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
			},
			ActiveHouseholdID:    helper.exampleHousehold.ID,
			HouseholdPermissions: helper.examplePermCheckers,
		}

		mockUserDataManager := &mocktypes.UserDataManagerMock{}
		mockUserDataManager.On(
			"GetSessionContextDataForUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(sessionCtxData, nil)
		helper.service.userDataManager = mockUserDataManager

		h := &testutils.MockHTTPHandler{}
		h.On(
			"ServeHTTP",
			testutils.HTTPResponseWriterMatcher,
			testutils.HTTPRequestMatcher,
		).Return()

		helper.req = helper.req.WithContext(context.WithValue(helper.ctx, sessions.SessionContextDataKey, sessionCtxData))

		mh := &testutils.MockHTTPHandler{}
		helper.service.AuthorizationMiddleware(mh).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusForbidden, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mh)
	})

	T.Run("with missing session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(*http.Request) (*sessions.ContextData, error) {
			return nil, nil
		}

		mh := &testutils.MockHTTPHandler{}
		helper.service.AuthorizationMiddleware(mh).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mh)
	})

	T.Run("without authorization for household", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
			},
			ActiveHouseholdID:    helper.exampleHousehold.ID,
			HouseholdPermissions: helper.examplePermCheckers,
		}

		sessionCtxData.HouseholdPermissions = map[string]authorization.HouseholdRolePermissionsChecker{}
		helper.service.sessionContextDataFetcher = func(*http.Request) (*sessions.ContextData, error) {
			return sessionCtxData, nil
		}

		mp := &mockmetrics.Int64Counter{}
		mp.On("Add", testutils.ContextMatcher, int64(1), mock.Anything).Return()
		helper.service.rejectedRequestCounter = mp

		helper.req = helper.req.WithContext(context.WithValue(helper.ctx, sessions.SessionContextDataKey, sessionCtxData))

		helper.service.AuthorizationMiddleware(&testutils.MockHTTPHandler{}).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		mock.AssertExpectationsForObjects(t, mp)
	})
}

func TestAuthenticationService_PermissionFilterMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.exampleUser.ServiceRole = authorization.ServiceAdminRole.String()
		helper.setContextFetcher(t)

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
			},
			ActiveHouseholdID:    helper.exampleHousehold.ID,
			HouseholdPermissions: helper.examplePermCheckers,
		}

		helper.req = helper.req.WithContext(context.WithValue(helper.req.Context(), sessions.SessionContextDataKey, sessionCtxData))

		mockHandler := &testutils.MockHTTPHandler{}
		mockHandler.On(
			"ServeHTTP",
			testutils.HTTPResponseWriterMatcher,
			testutils.HTTPRequestMatcher,
		).Return()

		helper.service.PermissionFilterMiddleware(authorization.InviteUserToHouseholdPermission)(mockHandler).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockHandler)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.exampleUser.ServiceRole = authorization.ServiceAdminRole.String()
		helper.setContextFetcher(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*sessions.ContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.PermissionFilterMiddleware(authorization.InviteUserToHouseholdPermission)(nil).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("unauthorized for household", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.exampleUser.ServiceRole = authorization.ServiceAdminRole.String()
		helper.setContextFetcher(t)

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(),
			},
			ActiveHouseholdID:    "different household, lol",
			HouseholdPermissions: helper.examplePermCheckers,
		}

		helper.req = helper.req.WithContext(context.WithValue(helper.req.Context(), sessions.SessionContextDataKey, sessionCtxData))
		helper.service.sessionContextDataFetcher = func(*http.Request) (*sessions.ContextData, error) {
			return sessionCtxData, nil
		}

		helper.service.PermissionFilterMiddleware(authorization.InviteUserToHouseholdPermission)(nil).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("without permission to perform action", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.exampleUser.ServiceRole = authorization.ServiceUserRole.String()
		helper.setContextFetcher(t)

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(authorization.InviteUserToHouseholdPermission.ID()),
			},
			ActiveHouseholdID: helper.exampleHousehold.ID,
			HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
				helper.exampleHousehold.ID: authorization.NewHouseholdRolePermissionChecker(authorization.InviteUserToHouseholdPermission.ID()),
			},
		}

		helper.req = helper.req.WithContext(context.WithValue(helper.req.Context(), sessions.SessionContextDataKey, sessionCtxData))
		helper.service.sessionContextDataFetcher = func(*http.Request) (*sessions.ContextData, error) {
			return sessionCtxData, nil
		}

		helper.service.PermissionFilterMiddleware(authorization.ArchiveHouseholdPermission)(nil).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})
}

func TestAuthenticationService_AdminMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.exampleUser.ServiceRole = authorization.ServiceAdminRole.String()
		helper.setContextFetcher(t)

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
			},
			ActiveHouseholdID:    helper.exampleHousehold.ID,
			HouseholdPermissions: helper.examplePermCheckers,
		}

		helper.req = helper.req.WithContext(context.WithValue(helper.req.Context(), sessions.SessionContextDataKey, sessionCtxData))

		mockHandler := &testutils.MockHTTPHandler{}
		mockHandler.On(
			"ServeHTTP",
			testutils.HTTPResponseWriterMatcher,
			testutils.HTTPRequestMatcher,
		).Return()

		helper.service.ServiceAdminMiddleware(mockHandler).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockHandler)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.exampleUser.ServiceRole = authorization.ServiceAdminRole.String()
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
			},
			ActiveHouseholdID:    helper.exampleHousehold.ID,
			HouseholdPermissions: helper.examplePermCheckers,
		}

		helper.req = helper.req.WithContext(context.WithValue(helper.req.Context(), sessions.SessionContextDataKey, sessionCtxData))

		mockHandler := &testutils.MockHTTPHandler{}
		helper.service.ServiceAdminMiddleware(mockHandler).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockHandler)
	})

	T.Run("with non-admin user", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		sessionCtxData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:                   helper.exampleUser.ID,
				AccountStatus:            helper.exampleUser.AccountStatus,
				AccountStatusExplanation: helper.exampleUser.AccountStatusExplanation,
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(helper.exampleUser.ServiceRole),
			},
			ActiveHouseholdID:    helper.exampleHousehold.ID,
			HouseholdPermissions: helper.examplePermCheckers,
		}

		helper.req = helper.req.WithContext(context.WithValue(helper.req.Context(), sessions.SessionContextDataKey, sessionCtxData))

		mockHandler := &testutils.MockHTTPHandler{}
		helper.service.ServiceAdminMiddleware(mockHandler).ServeHTTP(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockHandler)
	})
}
