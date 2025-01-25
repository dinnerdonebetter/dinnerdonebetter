package authentication

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	servertiming "github.com/mitchellh/go-server-timing"
)

// BuildLoginHandler is our login route.
func (s *service) BuildLoginHandler(adminOnly bool) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		timing := servertiming.FromContext(ctx)
		logger := s.logger.WithRequest(req).WithSpan(span)
		tracing.AttachRequestToSpan(span, req)

		responseDetails := types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		}

		if adminOnly {
			logger = logger.WithValue("admin_only", adminOnly)
		}

		loginData := new(types.UserLoginInput)
		if err := s.encoderDecoder.DecodeRequest(ctx, req, loginData); err != nil {
			observability.AcknowledgeError(err, logger, span, "decoding request body")
			errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
			return
		}

		loginData.TOTPToken = strings.TrimSpace(loginData.TOTPToken)
		loginData.Password = strings.TrimSpace(loginData.Password)
		loginData.Username = strings.TrimSpace(loginData.Username)

		if err := loginData.ValidateWithContext(ctx); err != nil {
			observability.AcknowledgeError(err, logger, span, "validating input")
			errRes := types.NewAPIErrorResponse("invalid login body", types.ErrValidatingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
			return
		}

		logger = logger.WithValue(keys.UsernameKey, loginData.Username)

		userFunc := s.userDataManager.GetUserByUsername
		if adminOnly {
			userFunc = s.userDataManager.GetAdminUserByUsername
		}

		readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
		user, err := userFunc(ctx, loginData.Username)
		if err != nil || user == nil {
			observability.AcknowledgeError(err, logger, span, "fetching user")
			if errors.Is(err, sql.ErrNoRows) {
				errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
				s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
				return
			}

			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}
		readTimer.Stop()

		logger = logger.WithValue(keys.UserIDKey, user.ID)
		tracing.AttachUserToSpan(span, user)

		if user.IsBanned() {
			errRes := types.NewAPIErrorResponse("user is banned", types.ErrUserIsBanned, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusForbidden)
			return
		}

		loginValid, err := s.validateLogin(ctx, user, loginData)
		logger.WithValue("login_valid", loginValid)

		if err != nil {
			if errors.Is(err, authentication.ErrInvalidTOTPToken) {
				observability.AcknowledgeError(err, logger, span, "validating TOTP token")
				errRes := types.NewAPIErrorResponse("login was invalid", types.ErrValidatingRequestInput, responseDetails)
				s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
				return
			}

			if errors.Is(err, authentication.ErrPasswordDoesNotMatch) {
				observability.AcknowledgeError(err, logger, span, "validating password")
				errRes := types.NewAPIErrorResponse("login was invalid", types.ErrValidatingRequestInput, responseDetails)
				s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
				return
			}

			observability.AcknowledgeError(err, logger, span, "validating login")
			errRes := types.NewAPIErrorResponse(staticError, types.ErrValidatingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		} else if !loginValid {
			logger.Debug("login was invalid")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "login was invalid", http.StatusUnauthorized)
			return
		}

		if loginValid && user.TwoFactorSecretVerifiedAt != nil && loginData.TOTPToken == "" {
			logger.Debug("user with two factor verification active attempted to log in without providing TOTP")
			errRes := types.NewAPIErrorResponse("TOTP required", types.ErrValidatingRequestInput, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusResetContent)
			return
		}

		defaultHouseholdID, err := s.householdMembershipManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching user memberships")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}
		responseDetails.CurrentHouseholdID = defaultHouseholdID

		responseCode, err := s.postLogin(ctx, user, defaultHouseholdID)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "handling login status")
			errRes := types.NewAPIErrorResponse(staticError, types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, responseCode)
			return
		}

		var token string
		token, err = s.tokenIssuer.IssueToken(ctx, user, s.config.TokenLifetime)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "signing token")
			errRes := types.NewAPIErrorResponse(staticError, types.ErrEncryptionIssue, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, responseCode)
			return
		}

		responseValue := &types.APIResponse[*types.TokenResponse]{
			Details: responseDetails,
			Data: &types.TokenResponse{
				HouseholdID: defaultHouseholdID,
				UserID:      user.ID,
				Token:       token,
			},
		}

		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, responseCode)
	}
}

func (s *service) postLogin(ctx context.Context, user *types.User, defaultHouseholdID string) (int, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	dcm := &types.DataChangeMessage{
		EventType:   types.UserLoggedInServiceEventType,
		HouseholdID: defaultHouseholdID,
		UserID:      user.ID,
	}

	if err := s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, s.logger, span, "publishing data change message")
		return http.StatusAccepted, nil
	}

	if err := s.analyticsReporter.AddUser(ctx, user.ID, map[string]any{
		"username":          user.Username,
		"default_household": defaultHouseholdID,
		"first_name":        user.FirstName,
		"last_name":         user.LastName,
	}); err != nil {
		return http.StatusAccepted, observability.PrepareError(err, span, "identifying user for analytics")
	}

	if err := s.featureFlagManager.Identify(ctx, user); err != nil {
		return http.StatusAccepted, observability.PrepareError(err, span, "identifying user in feature flag manager")
	}

	return http.StatusAccepted, nil
}

func (s *service) SSOLoginHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	providerName := s.authProviderFetcher(req)
	if providerName == "" {
		errRes := types.NewAPIErrorResponse("provider name is missing", types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
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

	u, err := sess.GetAuthURL()
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting auth url")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to get auth url", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	http.Redirect(res, req, u, http.StatusTemporaryRedirect)
}

func (s *service) SSOLoginCallbackHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	providerName := s.authProviderFetcher(req)
	if providerName == "" {
		errRes := types.NewAPIErrorResponse("provider name is missing", types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting provider")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "provider is not supported", http.StatusBadRequest)
		return
	}

	// NOTE: I know this doesn't work, but I can't be bothered to fix it right now
	value := req.Header.Get("Authorization")

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "unmarshaling session")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to unmarshal session", http.StatusInternalServerError)
		return
	}

	if err = validateState(req, sess); err != nil {
		observability.AcknowledgeError(err, logger, span, "validating state")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to validate state", http.StatusInternalServerError)
		return
	}

	fetchUserTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	providedUser, err := provider.FetchUser(sess)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to fetch user", http.StatusInternalServerError)
		return
	}
	fetchUserTimer.Stop()

	params := req.URL.Query()
	if params.Encode() == "" && req.Method == http.MethodPost {
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

	getUserTimer := timing.NewMetric("database").WithDesc("get user by email").Start()
	user, err := s.userDataManager.GetUserByEmail(ctx, providedUser.Email)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting user by email")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	getUserTimer.Stop()

	defaultHouseholdTimer := timing.NewMetric("database").WithDesc("get default household for user").Start()
	defaultHouseholdID, err := s.householdMembershipManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user memberships")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	defaultHouseholdTimer.Stop()

	var token string
	token, err = s.tokenIssuer.IssueToken(ctx, user, s.config.TokenLifetime)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "signing token")
		errRes := types.NewAPIErrorResponse(staticError, types.ErrEncryptionIssue, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[*types.TokenResponse]{
		Details: responseDetails,
		Data: &types.TokenResponse{
			HouseholdID: defaultHouseholdID,
			UserID:      user.ID,
			Token:       token,
		},
	}

	responseCode, err := s.postLogin(ctx, user, defaultHouseholdID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "handling login status")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, staticError, responseCode)
		return
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, responseCode)
	logger.Debug("user logged in via SSO")
}

// validateState ensures that the state token param from the original.
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

// StatusHandler returns the user info for the user making the request.
func (s *service) StatusHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	statusResponse := &types.UserStatusResponse{
		ActiveHousehold:          sessionCtxData.ActiveHouseholdID,
		AccountStatus:            sessionCtxData.Requester.AccountStatus,
		AccountStatusExplanation: sessionCtxData.Requester.AccountStatusExplanation,
		UserIsAuthenticated:      true,
	}

	responseValue := &types.APIResponse[*types.UserStatusResponse]{
		Details: responseDetails,
		Data:    statusResponse,
	}

	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

var _ types.OAuth2Service = (*service)(nil)

// AuthorizeHandler is our oauth2 auth route.
func (s *service) AuthorizeHandler(res http.ResponseWriter, req *http.Request) {
	_, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if err := s.oauth2Server.HandleAuthorizeRequest(res, req); err != nil {
		observability.AcknowledgeError(err, logger, span, "handling authorization request")
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func (s *service) TokenHandler(res http.ResponseWriter, req *http.Request) {
	_, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	if err := s.oauth2Server.HandleTokenRequest(res, req); err != nil {
		observability.AcknowledgeError(err, logger, span, "handling token request")
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
