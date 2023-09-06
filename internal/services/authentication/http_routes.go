package authentication

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/gorilla/securecookie"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

var (
	customCookieDomainHeader = "X-DDB-COOKIE-DOMAIN"

	allowedCookiesHat    sync.Mutex
	allowedCookieDomains = map[string]uint{
		".dinnerdonebetter.local": 0,
		".dinnerdonebetter.dev":   1,
		".dinnerdonebetter.com":   2,
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

		userFunc := s.userDataManager.GetUserByUsername
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

		responseCode, err := s.postLogin(ctx, user, defaultHouseholdID, req, res)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "handling login status")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, responseCode)
			return
		}

		statusResponse := &types.UserStatusResponse{
			UserID:                   user.ID,
			UserIsAuthenticated:      true,
			AccountStatus:            user.AccountStatus,
			ActiveHousehold:          defaultHouseholdID,
			AccountStatusExplanation: user.AccountStatusExplanation,
		}

		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, statusResponse, responseCode)
		logger.Debug("user logged in")
	}
}

func (s *service) postLogin(ctx context.Context, user *types.User, defaultHouseholdID string, req *http.Request, res http.ResponseWriter) (int, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	requestedCookieDomain := s.determineCookieDomain(ctx, req)

	cookie, err := s.issueSessionManagedCookie(ctx, defaultHouseholdID, user.ID, requestedCookieDomain)
	if err != nil {
		return http.StatusInternalServerError, observability.PrepareError(err, span, "issuing cookie")
	}

	http.SetCookie(res, cookie)

	dcm := &types.DataChangeMessage{
		EventType:   types.UserLoggedInCustomerEventType,
		HouseholdID: defaultHouseholdID,
		UserID:      user.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return http.StatusAccepted, observability.PrepareError(err, span, "publishing data change message")
	}

	if err = s.analyticsReporter.AddUser(ctx, user.ID, map[string]any{
		"username":          user.Username,
		"default_household": defaultHouseholdID,
		"first_name":        user.FirstName,
		"last_name":         user.LastName,
	}); err != nil {
		return http.StatusAccepted, observability.PrepareError(err, span, "identifying user for analytics")
	}

	if err = s.featureFlagManager.Identify(ctx, user); err != nil {
		return http.StatusAccepted, observability.PrepareError(err, span, "identifying user in feature flag manager")
	}

	return http.StatusAccepted, nil
}

func (s *service) SSOProviderHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	providerName := s.authProviderFetcher(req)
	if providerName == "" {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "provider name is missing", http.StatusBadRequest)
		return
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting provider")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "provider is not supported", http.StatusBadRequest)
		return
	}

	sess, err := provider.BeginAuth(gothic.SetState(req))
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "beginning auth")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to begin auth", http.StatusInternalServerError)
		return
	}

	ctx, err = s.sessionManager.Load(ctx, providerName)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "loading session")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to load session", http.StatusInternalServerError)
		return
	}

	s.sessionManager.Put(ctx, providerName, sess.Marshal())

	u, err := sess.GetAuthURL()
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting auth url")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to get auth url", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	http.Redirect(res, req, u, http.StatusTemporaryRedirect)
}

func (s *service) SSOCallbackHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	providerName := s.authProviderFetcher(req)
	if providerName == "" {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "provider name is missing", http.StatusBadRequest)
		return
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting provider")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "provider is not supported", http.StatusBadRequest)
		return
	}

	rawValue := s.sessionManager.Get(ctx, providerName)
	value, ok := rawValue.(string)
	if !ok {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to get session", http.StatusInternalServerError)
		return
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "unmarshalling session")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to unmarshal session", http.StatusInternalServerError)
		return
	}

	err = validateState(req, sess)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "validating state")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to validate state", http.StatusInternalServerError)
		return
	}

	providedUser, err := provider.FetchUser(sess)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to fetch user", http.StatusInternalServerError)
		return
	}

	params := req.URL.Query()
	if params.Encode() == "" && req.Method == "POST" {
		if err = req.ParseForm(); err != nil {
			observability.AcknowledgeError(err, logger, span, "parsing form")
		}
		params = req.Form
	}

	// get new token and retry fetch
	if _, err = sess.Authorize(provider, params); err != nil {
		observability.AcknowledgeError(err, logger, span, "authorizing session")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to authorize session", http.StatusInternalServerError)
		return
	}

	user, err := s.userDataManager.GetUserByEmail(ctx, providedUser.Email)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting user by email")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to get user by email", http.StatusInternalServerError)
		return
	}

	defaultHouseholdID, err := s.householdMembershipManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user memberships")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, http.StatusInternalServerError)
		return
	}

	responseCode, err := s.postLogin(ctx, user, defaultHouseholdID, req, res)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "handling login status")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, responseCode)
		return
	}

	statusResponse := &types.UserStatusResponse{
		UserID:                   user.ID,
		UserIsAuthenticated:      true,
		AccountStatus:            user.AccountStatus,
		ActiveHousehold:          defaultHouseholdID,
		AccountStatusExplanation: user.AccountStatusExplanation,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, statusResponse, responseCode)
	logger.Debug("user logged in via SSO")
}

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateState(req *http.Request, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	reqState := gothic.GetState(req)

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != reqState) {
		return errors.New("state token mismatch")
	}
	return nil
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

	dcm := &types.DataChangeMessage{
		EventType: types.UserChangedActiveHouseholdCustomerEventType,
		Context: map[string]any{
			"old_household_id": sessionCtxData.ActiveHouseholdID,
		},
		HouseholdID: householdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
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

	dcm := &types.DataChangeMessage{
		EventType: types.UserLoggedOutCustomerEventType,
		UserID:    sessionCtxData.Requester.UserID,
	}

	if dataPublishErr := s.dataChangesPublisher.Publish(ctx, dcm); dataPublishErr != nil {
		observability.AcknowledgeError(dataPublishErr, logger, span, "publishing data change message")
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

var _ types.OAuth2Service = (*service)(nil)

// AuthorizeHandler is our oauth2 auth route.
func (s *service) AuthorizeHandler(res http.ResponseWriter, req *http.Request) {
	if err := s.oauth2Server.HandleAuthorizeRequest(res, req); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func (s *service) TokenHandler(res http.ResponseWriter, req *http.Request) {
	if err := s.oauth2Server.HandleTokenRequest(res, req); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
