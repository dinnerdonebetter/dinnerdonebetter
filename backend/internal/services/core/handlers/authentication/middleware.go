package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	authentication2 "github.com/dinnerdonebetter/backend/internal/lib/authentication/sessioncontext"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	oauth2errors "github.com/go-oauth2/oauth2/v4/errors"
	servertiming "github.com/mitchellh/go-server-timing"
)

const (
	zuckModeUserHeader      = "X-DDB-Zuck-Mode-User"
	zuckModeHouseholdHeader = "X-DDB-Zuck-Mode-Household"
)

var (
	// ErrUserNotAuthorizedToImpersonateOthers is returned when a user is not authorized to impersonate others.
	ErrUserNotAuthorizedToImpersonateOthers = errors.New("user not authorized to impersonate others")
)

func (s *service) determineZuckMode(ctx context.Context, req *http.Request, sessionContextData *authentication2.SessionContextData) (userID, householdID string, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req).WithSpan(span)

	if zuckUserID := req.Header.Get(zuckModeUserHeader); zuckUserID != "" {
		if !sessionContextData.ServiceRolePermissionChecker().CanImpersonateUsers() {
			return "", "", ErrUserNotAuthorizedToImpersonateOthers
		}

		if _, err = s.userDataManager.GetUser(ctx, zuckUserID); err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching user info for zuck mode")
			return "", "", err
		}

		if zuckHouseholdID := req.Header.Get(zuckModeHouseholdHeader); zuckHouseholdID == "" {
			householdID, err = s.householdMembershipManager.GetDefaultHouseholdIDForUser(ctx, zuckUserID)
			if err != nil {
				observability.AcknowledgeError(err, logger, span, "fetching household info for zuck mode")
				return "", "", err
			}
		} else {
			return zuckUserID, zuckHouseholdID, nil
		}

		return zuckUserID, householdID, nil
	}

	return "", "", nil
}

// UserAttributionMiddleware is concerned with figuring out who a user is, but not worried about kicking out users who are not known.
func (s *service) UserAttributionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		timing := servertiming.FromContext(ctx)
		logger := s.logger.WithRequest(req).WithSpan(span)
		responseDetails := types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		}

		// validate bearer token.
		tokenTimer := timing.NewMetric("token validation").WithDesc("validating bearer token from request").Start()
		token, err := s.oauth2Server.ValidationBearerToken(req)
		if err != nil {
			if errors.Is(err, oauth2errors.ErrExpiredAccessToken) || errors.Is(err, oauth2errors.ErrExpiredRefreshToken) {
				res.WriteHeader(http.StatusTeapot)
				return
			} else {
				logger.Error("determining user ID", err)
			}
		}
		tokenTimer.Stop()

		if token != nil {
			if userID := token.GetUserID(); userID != "" {
				userAttributionTimer := timing.NewMetric("user attribution").WithDesc("attributing user to request").Start()
				sessionCtxData, sessionCtxDataErr := s.householdMembershipManager.BuildSessionContextDataForUser(ctx, userID)
				if sessionCtxDataErr != nil {
					observability.AcknowledgeError(sessionCtxDataErr, logger, span, "fetching user info for cookie")
					errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
					s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
					return
				}

				zuckUserID, zuckHouseholdID, zuckErr := s.determineZuckMode(ctx, req, sessionCtxData)
				if zuckErr != nil {
					observability.AcknowledgeError(zuckErr, logger, span, "fetching user info for zuck mode")
					errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
					s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
					return
				}

				if zuckUserID != "" {
					sessionCtxData.Requester.UserID = zuckUserID
				}

				if zuckHouseholdID != "" {
					sessionCtxData.ActiveHouseholdID = zuckHouseholdID
					sessionCtxData.HouseholdPermissions[zuckHouseholdID] = authorization.NewHouseholdRolePermissionChecker(authorization.HouseholdMemberRole.String())
				}

				userAttributionTimer.Stop()
				if sessionCtxData != nil {
					next.ServeHTTP(res, req.WithContext(context.WithValue(ctx, authentication2.SessionContextDataKey, sessionCtxData)))
					return
				}
			}
		}

		next.ServeHTTP(res, req)
	})
}

// AuthorizationMiddleware checks to see if a user is associated with the request, and then determines whether said request can proceed.
func (s *service) AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req).WithSpan(span)

		// UserAttributionMiddleware should be called before this middleware.
		if sessionCtxData, err := s.sessionContextDataFetcher(req); err == nil && sessionCtxData != nil {
			tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
			logger = sessionCtxData.AttachToLogger(logger)

			if sessionCtxData.Requester.AccountStatus == string(types.BannedUserAccountStatus) || sessionCtxData.Requester.AccountStatus == string(types.TerminatedUserAccountStatus) {
				logger.Info("banned user attempted to make request")
				s.rejectedRequestCounter.Add(ctx, 1)
				http.Redirect(res, req, "/", http.StatusForbidden)
				return
			}

			canImpersonateUsers := sessionCtxData.ServiceRolePermissionChecker().CanImpersonateUsers()

			if _, authorizedForHousehold := sessionCtxData.HouseholdPermissions[sessionCtxData.ActiveHouseholdID]; !authorizedForHousehold && !canImpersonateUsers {
				logger.Info("user trying to access household they are not authorized for")
				s.rejectedRequestCounter.Add(ctx, 1)
				http.Redirect(res, req, "/", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(res, req)
			return
		}

		logger.Info("no user attached to request")
		http.Redirect(res, req, "/users/login", http.StatusUnauthorized)
	})
}

// PermissionFilterMiddleware filters users out of requests based on provided functions.
func (s *service) PermissionFilterMiddleware(permissions ...authorization.Permission) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx, span := s.tracer.StartSpan(req.Context())
			defer span.End()

			timing := servertiming.FromContext(ctx)
			logger := s.logger.WithRequest(req).WithSpan(span)
			logger.Debug("checking permissions in middleware")

			unauthorizedResponse := &types.APIResponse[any]{
				Details: types.ResponseDetails{
					TraceID: span.SpanContext().TraceID().String(),
				},
				Error: &types.APIError{
					Message: "invalid credentials provided",
				},
			}

			// check for a session context data first.
			sessionContextData, err := s.sessionContextDataFetcher(req)
			if err != nil {
				observability.AcknowledgeError(err, logger, span, "retrieving session context data")
				s.encoderDecoder.EncodeResponseWithStatus(ctx, res, unauthorizedResponse, http.StatusUnauthorized)
				return
			}

			permissionCheckTimer := timing.NewMetric("permissions check").WithDesc("checking user permissions").Start()
			logger = sessionContextData.AttachToLogger(logger)
			isServiceAdmin := sessionContextData.Requester.ServicePermissions.IsServiceAdmin()
			logger = logger.WithValue("is_service_admin", isServiceAdmin)

			if _, allowed := sessionContextData.HouseholdPermissions[sessionContextData.ActiveHouseholdID]; !allowed && !isServiceAdmin {
				permissionCheckTimer.Stop()
				logger.Info("not authorized for household")
				s.encoderDecoder.EncodeResponseWithStatus(ctx, res, unauthorizedResponse, http.StatusUnauthorized)
				return
			}

			for _, perm := range permissions {
				doesNotHaveServicePermission := !sessionContextData.ServiceRolePermissionChecker().HasPermission(perm)
				doesNotHaveHouseholdPermission := !sessionContextData.HouseholdRolePermissionsChecker().HasPermission(perm)
				if doesNotHaveServicePermission && doesNotHaveHouseholdPermission {
					permissionCheckTimer.Stop()
					logger.WithValue("deficient_permission", perm.ID()).Info("request filtered out")
					s.encoderDecoder.EncodeResponseWithStatus(ctx, res, unauthorizedResponse, http.StatusUnauthorized)
					return
				}
			}

			permissionCheckTimer.Stop()
			next.ServeHTTP(res, req)
		})
	}
}

// ServiceAdminMiddleware restricts requests to admin users only.
func (s *service) ServiceAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req).WithSpan(span)

		responseDetails := types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		}

		sessionCtxData, err := s.sessionContextDataFetcher(req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "retrieving session context data")
			errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
			return
		}

		tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
		logger = sessionCtxData.AttachToLogger(logger)

		if !sessionCtxData.Requester.ServicePermissions.IsServiceAdmin() {
			logger.Debug("ServiceAdminMiddleware called by non-admin user")
			errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrUserIsNotAuthorized, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}
