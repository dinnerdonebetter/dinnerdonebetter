package authentication

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/o1egl/paseto"

	"github.com/prixfixeco/backend/internal/authentication"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

var (
	customCookieDomainHeader = "X-PRIXFIXE-COOKIE-DOMAIN"

	allowedCookiesHat    sync.Mutex
	allowedCookieDomains = map[string]uint{
		".prixfixe.local": 0,
		".prixfixe.dev":   1,
		".prixfixe.app":   2,
	}
)

// determineCookieDomain determines which domain to assign a cookie.
func (s *service) determineCookieDomain(ctx context.Context, req *http.Request) string {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	requestedCookieDomain := s.config.Cookies.Domain
	if headerCookieDomain := req.Header.Get(customCookieDomainHeader); headerCookieDomain != "" {
		allowedCookiesHat.Lock()
		// if the requested domain is present in the map, and it has a lower score than the current domain, then
		if currentScore, ok1 := allowedCookieDomains[requestedCookieDomain]; ok1 {
			if newScore, ok2 := allowedCookieDomains[headerCookieDomain]; ok2 {
				if currentScore > newScore {
					requestedCookieDomain = headerCookieDomain
				}
			}
		}
		allowedCookiesHat.Unlock()
	}

	return requestedCookieDomain
}

// BuildLoginHandler is our login route.
func (s *service) BuildLoginHandler(adminOnly bool) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)
		tracing.AttachRequestToSpan(span, req)

		if adminOnly {
			logger = logger.WithValue("admin_only", adminOnly)
		}

		loginData := new(types.UserLoginInput)
		if err := s.encoderDecoder.DecodeRequest(ctx, req, loginData); err != nil {
			observability.AcknowledgeError(err, logger, span, "decoding request body")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
			return
		}

		loginData.TOTPToken = strings.TrimSpace(loginData.TOTPToken)
		loginData.Password = strings.TrimSpace(loginData.Password)
		loginData.Username = strings.TrimSpace(loginData.Username)

		if err := loginData.ValidateWithContext(ctx, s.config.MinimumUsernameLength, s.config.MinimumPasswordLength); err != nil {
			observability.AcknowledgeError(err, logger, span, "validating input")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid login body", http.StatusBadRequest)
			return
		}

		logger = logger.WithValue(keys.UsernameKey, loginData.Username)

		requestedCookieDomain := s.determineCookieDomain(ctx, req)
		if requestedCookieDomain != "" {
			logger = logger.WithValue("cookie_domain", requestedCookieDomain)
		}

		var userFunc = s.userDataManager.GetUserByUsername
		if adminOnly {
			userFunc = s.userDataManager.GetAdminUserByUsername
		}

		user, err := userFunc(ctx, loginData.Username)
		if err != nil || user == nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
				return
			}

			observability.AcknowledgeError(err, logger, span, "fetching user")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
			return
		}

		logger = logger.WithValue(keys.UserIDKey, user.ID)
		tracing.AttachUserToSpan(span, user)

		if user.IsBanned() {
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "user is banned", http.StatusForbidden)
			return
		}

		loginValid, err := s.validateLogin(ctx, user, loginData)
		logger.WithValue("login_valid", loginValid)

		if err != nil {
			if errors.Is(err, authentication.ErrInvalidTOTPToken) {
				observability.AcknowledgeError(err, logger, span, "validating TOTP token")
				s.encoderDecoder.EncodeErrorResponse(ctx, res, "login was invalid", http.StatusUnauthorized)
				return
			}

			if errors.Is(err, authentication.ErrPasswordDoesNotMatch) {
				observability.AcknowledgeError(err, logger, span, "validating password")
				s.encoderDecoder.EncodeErrorResponse(ctx, res, "login was invalid", http.StatusUnauthorized)
				return
			}

			observability.AcknowledgeError(err, logger, span, "validating login")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
			return
		} else if !loginValid {
			logger.Debug("login was invalid")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "login was invalid", http.StatusUnauthorized)
			return
		}

		if loginValid && user.TwoFactorSecretVerifiedAt != nil && loginData.TOTPToken == "" {
			logger.Debug("user with two factor verification active attempted to log in without providing TOTP")
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, "TOTP required", http.StatusResetContent)
			return
		}

		defaultHouseholdID, err := s.householdMembershipManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching user memberships")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
			return
		}

		cookie, err := s.issueSessionManagedCookie(ctx, defaultHouseholdID, user.ID, requestedCookieDomain)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "issuing cookie")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
			return
		}

		if s.dataChangesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:             types.UserDataType,
				EventType:            types.UserLoggedInCustomerEventType,
				HouseholdID:          defaultHouseholdID,
				AttributableToUserID: user.ID,
			}

			if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
				observability.AcknowledgeError(err, logger, span, "publishing data change message")
			}
		}

		http.SetCookie(res, cookie)

		statusResponse := &types.UserStatusResponse{
			UserID:                   user.ID,
			UserIsAuthenticated:      true,
			AccountStatus:            user.AccountStatus,
			ActiveHousehold:          defaultHouseholdID,
			AccountStatusExplanation: user.AccountStatusExplanation,
		}

		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, statusResponse, http.StatusAccepted)
		logger.Debug("user logged in")
	}
}

// ChangeActiveHouseholdHandler is our login route.
func (s *service) ChangeActiveHouseholdHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	input := new(types.ChangeActiveHouseholdInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	householdID := input.HouseholdID
	logger = logger.WithValue("new_session_household_id", householdID)

	requesterID := sessionCtxData.Requester.UserID
	logger = logger.WithValue("user_id", requesterID)

	authorizedForHousehold, err := s.householdMembershipManager.UserIsMemberOfHousehold(ctx, requesterID, householdID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "checking permissions")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
		return
	}

	if !authorizedForHousehold {
		logger.Debug("invalid household ID requested for activation")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	requestedCookieDomain := s.determineCookieDomain(ctx, req)
	if requestedCookieDomain != "" {
		logger = logger.WithValue("cookie_domain", requestedCookieDomain)
	}

	cookie, err := s.issueSessionManagedCookie(ctx, householdID, requesterID, requestedCookieDomain)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "issuing cookie")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.UserDataType,
			EventType:            types.UserChangedActiveHouseholdCustomerEventType,
			AttributableToUserID: requesterID,
			Context: map[string]string{
				"old_household_id": sessionCtxData.ActiveHouseholdID,
			},
			HouseholdID: householdID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	logger.Info("successfully changed active session household")
	http.SetCookie(res, cookie)

	res.WriteHeader(http.StatusAccepted)
}

// EndSessionHandler is our logout route.
func (s *service) EndSessionHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, fetchSessionContextErr := s.sessionContextDataFetcher(req)
	if fetchSessionContextErr != nil {
		observability.AcknowledgeError(fetchSessionContextErr, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	ctx, loadErr := s.sessionManager.Load(ctx, "")
	if loadErr != nil {
		// this can literally never happen in this version of scs, because the token is empty
		observability.AcknowledgeError(loadErr, logger, span, "loading token")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if destroyErr := s.sessionManager.Destroy(ctx); destroyErr != nil {
		observability.AcknowledgeError(destroyErr, logger, span, "destroying session")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	requestedCookieDomain := s.determineCookieDomain(ctx, req)
	if requestedCookieDomain != "" {
		logger = logger.WithValue("cookie_domain", requestedCookieDomain)
	}

	newCookie, cookieBuildingErr := s.buildLogoutCookie(ctx, req)
	if cookieBuildingErr != nil || newCookie == nil {
		observability.AcknowledgeError(cookieBuildingErr, logger, span, "building cookie")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	http.SetCookie(res, newCookie)

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.UserDataType,
			EventType:            types.UserLoggedOutCustomerEventType,
			AttributableToUserID: sessionCtxData.Requester.UserID,
		}

		if dataPublishErr := s.dataChangesPublisher.Publish(ctx, dcm); dataPublishErr != nil {
			observability.AcknowledgeError(dataPublishErr, logger, span, "publishing data change message")
		}
	}

	logger.Debug("user logged out")

	res.WriteHeader(http.StatusAccepted)
}

// StatusHandler returns the user info for the user making the request.
func (s *service) StatusHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	var statusResponse *types.UserStatusResponse

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	statusResponse = &types.UserStatusResponse{
		ActiveHousehold:          sessionCtxData.ActiveHouseholdID,
		AccountStatus:            sessionCtxData.Requester.AccountStatus,
		AccountStatusExplanation: sessionCtxData.Requester.AccountStatusExplanation,
		UserIsAuthenticated:      true,
	}

	s.encoderDecoder.RespondWithData(ctx, res, statusResponse)
}

const (
	pasetoRequestTimeThreshold = 2 * time.Minute
)

// PASETOHandler returns the user info for the user making the request.
func (s *service) PASETOHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.PASETOCreationInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	requestedHousehold := input.HouseholdID
	logger = logger.WithValue(keys.APIClientClientIDKey, input.ClientID)

	if requestedHousehold != "" {
		logger = logger.WithValue("requested_household", requestedHousehold)
	}

	reqTime := time.Unix(0, input.RequestTime)
	if time.Until(reqTime) > pasetoRequestTimeThreshold || time.Since(reqTime) > pasetoRequestTimeThreshold {
		logger.WithValue("provided_request_time", reqTime.String()).Debug("PASETO request denied because its time is out of threshold")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	sum, err := base64.RawURLEncoding.DecodeString(req.Header.Get(signatureHeaderKey))
	if err != nil || len(sum) == 0 {
		logger.WithValue("sum_length", len(sum)).Error(err, "invalid signature")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	client, clientRetrievalErr := s.apiClientManager.GetAPIClientByClientID(ctx, input.ClientID)
	if clientRetrievalErr != nil {
		observability.AcknowledgeError(err, logger, span, "fetching API client")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	mac := hmac.New(sha256.New, client.ClientSecret)
	if _, macWriteErr := mac.Write(s.encoderDecoder.MustEncodeJSON(ctx, input)); macWriteErr != nil {
		// sha256.digest.Write does not ever return an error, so this branch will remain "uncovered" :(
		observability.AcknowledgeError(err, logger, span, "writing HMAC message for comparison")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	if !hmac.Equal(sum, mac.Sum(nil)) {
		logger.Info("invalid credentials passed to PASETO creation route")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	user, err := s.userDataManager.GetUser(ctx, client.BelongsToUser)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving user")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	logger = logger.WithValue(keys.UserIDKey, user.ID)

	sessionCtxData, err := s.householdMembershipManager.BuildSessionContextDataForUser(ctx, user.ID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving perms for API client")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	var requestedHouseholdID string

	if requestedHousehold != "" {
		if _, isMember := sessionCtxData.HouseholdPermissions[requestedHousehold]; !isMember {
			logger.Debug("invalid household ID requested for token")
			s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
			return
		}

		logger.WithValue("requested_household", requestedHousehold).Debug("setting token household ID to requested household")
		requestedHouseholdID = requestedHousehold
		sessionCtxData.ActiveHouseholdID = requestedHousehold
	} else {
		requestedHouseholdID = sessionCtxData.ActiveHouseholdID
	}

	logger = logger.WithValue(keys.HouseholdIDKey, requestedHouseholdID)

	// Encrypt data
	tokenRes, err := s.buildPASETOResponse(ctx, sessionCtxData, client)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "encrypting PASETO")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger.Info("PASETO issued")

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, tokenRes, http.StatusAccepted)
}

func (s *service) buildPASETOToken(ctx context.Context, sessionCtxData *types.SessionContextData, client *types.APIClient) paseto.JSONToken {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	now := time.Now().UTC()
	lifetime := time.Duration(math.Min(float64(maxPASETOLifetime), float64(s.config.PASETO.Lifetime)))
	expiry := now.Add(lifetime)

	jsonToken := paseto.JSONToken{
		Audience:   client.BelongsToUser,
		Subject:    client.BelongsToUser,
		Jti:        uuid.NewString(),
		Issuer:     s.config.PASETO.Issuer,
		IssuedAt:   now,
		NotBefore:  now,
		Expiration: expiry,
	}

	jsonToken.Set(pasetoDataKey, base64.RawURLEncoding.EncodeToString(sessionCtxData.ToBytes()))

	return jsonToken
}

func (s *service) buildPASETOResponse(ctx context.Context, sessionCtxData *types.SessionContextData, client *types.APIClient) (*types.PASETOResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	jsonToken := s.buildPASETOToken(ctx, sessionCtxData, client)

	// Encrypt data
	token, err := paseto.NewV2().Encrypt(s.config.PASETO.LocalModeKey, jsonToken, "")
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting PASETO")
	}

	tokenRes := &types.PASETOResponse{
		Token:     token,
		ExpiresAt: jsonToken.Expiration.String(),
	}

	return tokenRes, nil
}

// CycleCookieSecretHandler rotates the cookie building secret with a new random secret.
func (s *service) CycleCookieSecretHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Info("cycling cookie secret!")

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	if !sessionCtxData.Requester.ServicePermissions.CanCycleCookieSecrets() {
		logger.Debug("invalid permissions")
		s.encoderDecoder.EncodeInvalidPermissionsResponse(ctx, res)
		return
	}

	s.cookieManager = securecookie.New(
		securecookie.GenerateRandomKey(cookieSecretSize),
		[]byte(s.config.Cookies.BlockKey),
	)

	res.WriteHeader(http.StatusAccepted)
}
