package integration

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityconverters "github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	identityfakes "github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	types "github.com/dinnerdonebetter/backend/internal/domain/oauth"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	"github.com/dinnerdonebetter/backend/internal/services/identity/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpTestServerAddress = "http://localhost:8000"
	grpcTestServerAddress = ":8001"

	adminUserPassword = "integration-tests-are-cool"

	nonexistentID = "00000000000000000000"
)

var (
	premadeAdminUser = &identity.User{
		ID:              identifiers.New(),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@example.email",
		Username:        "admin_user",
		HashedPassword:  adminUserPassword,
	}

	adminClient client.Client
)

func buildUnauthenticatedGRPCClientForTest(t *testing.T) client.Client {
	t.Helper()

	c, err := buildUnauthenticatedGRPCClient()
	require.NoError(t, err)

	return c
}

func buildUnauthenticatedGRPCClient() (client.Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return client.BuildClient(grpcTestServerAddress, opts...)
}

func buildAuthedGRPCClient(ctx context.Context, scopes []string, token string) client.Client {
	state, err := random.GenerateBase64EncodedString(ctx, 32)
	if err != nil {
		panic(err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     createdClientID,
		ClientSecret: createdClientSecret,
		Scopes:       scopes, // TODO: This should be nil-able
		RedirectURL:  httpTestServerAddress,
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   httpTestServerAddress + "/oauth2/authorize",
			TokenURL:  httpTestServerAddress + "/oauth2/token",
		},
	}

	authCodeURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		authCodeURL,
		http.NoBody,
	)
	if err != nil {
		panic(fmt.Errorf("failed to get oauth2 code: %w", err))
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Location", "localhost")

	httpClient := tracing.BuildTracedHTTPClient()
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := httpClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("failed to get oauth2 code: %w", err))
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println("failed to close oauth2 response body", err)
		}
	}()

	const (
		codeKey = "code"
	)

	rl, err := res.Location()
	if err != nil {
		panic(err)
	}

	code := rl.Query().Get(codeKey)
	if code == "" {
		panic("code not returned from oauth2 redirect")
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&insecureOAuth{
			TokenSource: oauth2Config.TokenSource(ctx, oauth2Token),
		}),
	}

	c, err := client.BuildClient(grpcTestServerAddress, opts...)
	if err != nil {
		panic(err)
	}

	return c
}

// Custom insecure OAuth2 credentials that skip security checks
type insecureOAuth struct {
	TokenSource oauth2.TokenSource
}

func (i *insecureOAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	token, err := i.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	return map[string]string{"authorization": token.Type() + " " + token.AccessToken}, nil
}

func (i *insecureOAuth) RequireTransportSecurity() bool {
	return false // Explicitly allow insecure transport
}

func deriveServerConfig() (*config.APIServiceConfig, error) {
	wd, _ := os.Getwd()
	fmt.Println(wd)

	content, err := os.ReadFile(apiConfigurationFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read api configuration file: %w", err)
	}

	decoder := encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	var x *config.APIServiceConfig
	if err = decoder.DecodeBytes(context.Background(), content, &x); err != nil {
		return nil, fmt.Errorf("failed to decode api configuration file: %w", err)
	}

	return x, nil
}

func createOAuth2ClientForTests(ctx context.Context, pgc database.Client, dbCfg *databasecfg.Config) error {
	auditRepo := auditlogentries.ProvideAuditLogRepository(nil, nil, pgc)
	oauth2ClientManager := oauth.ProvideOAuthRepository(nil, nil, auditRepo, *dbCfg, pgc)

	clientID, err := random.GenerateHexEncodedString(ctx, 16)
	if err != nil {
		return fmt.Errorf("failed to generate client ID: %w", err)
	}

	clientSecret, err := random.GenerateHexEncodedString(ctx, 16)
	if err != nil {
		return fmt.Errorf("failed to generate client secret: %w", err)
	}

	createdClient, err := oauth2ClientManager.CreateOAuth2Client(ctx, &types.OAuth2ClientDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "integration_client",
		Description:  "integration test client",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return fmt.Errorf("failed to create oauth2 client: %w", err)
	}

	createdClientID, createdClientSecret = createdClient.ClientID, createdClient.ClientSecret

	return nil
}

func createPremadeAdminUser(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, identityRepo identity.Repository, dbClient database.Client) (*identity.User, error) {
	hasher := authentication.ProvideArgon2Authenticator(logger, tracerProvider)

	actuallyHashedPass, err := hasher.HashPassword(ctx, premadeAdminUser.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	premadeAdminUser.HashedPassword = actuallyHashedPass

	var user *identity.User
	if user, err = identityRepo.GetUserByUsername(ctx, premadeAdminUser.Username); err == nil {
		return user, nil
	}

	user, err = identityRepo.CreateUser(ctx, identityconverters.ConvertUserToUserDatabaseCreationInput(premadeAdminUser))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// one-off query because I really don't want to make this functionality concrete
	if _, err = dbClient.DB().Exec(fmt.Sprintf("UPDATE users SET service_role='service_admin' WHERE id='%s'", user.ID)); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if err = identityRepo.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to mark user as verified: %w", err)
	}

	adminUser, err := identityRepo.GetAdminUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin user: %w", err)
	}

	return adminUser, nil
}

func hashStringToNumber(s string) uint64 {
	// Create a new FNV-1a 64-bit hash object
	h := fnv.New64a()

	// Write the bytes of the string into the hash object
	_, err := h.Write([]byte(s))
	if err != nil {
		// Handle error if necessary
		panic(err)
	}

	// Return the resulting hash value as a number (uint64)
	return h.Sum64()
}

func createServiceUserForTest(t *testing.T, verifyTOTP bool, in *identity.UserRegistrationInput) *identity.User {
	t.Helper()

	user, err := createServiceUser(t.Context(), verifyTOTP, in)
	require.NoError(t, err)

	return user
}

func createServiceUser(ctx context.Context, verifyTOTP bool, in *identity.UserRegistrationInput) (*identity.User, error) {
	c, err := buildUnauthenticatedGRPCClient()
	if err != nil {
		return nil, fmt.Errorf("initializing client: %w", err)
	}

	if in == nil {
		in = identityfakes.BuildFakeUserCreationInput()
	}
	input := converters.ConvertUserRegistrationInputToGRPCUserRegistrationInput(in)

	res, err := c.CreateUser(ctx, &identitysvc.CreateUserRequest{Input: input})
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	ucr := res.Created

	if verifyTOTP {
		if err = verifyTOTPSecretForUser(ctx, c, ucr.CreatedUserID, ucr.TwoFactorSecret); err != nil {
			return nil, fmt.Errorf("verifying totp code: %w", err)
		}
	}
	u := &identity.User{
		ID:              ucr.CreatedUserID,
		Username:        ucr.Username,
		EmailAddress:    ucr.EmailAddress,
		TwoFactorSecret: ucr.TwoFactorSecret,
		CreatedAt:       grpcconverters.ConvertPBTimestampToTime(ucr.CreatedAt),
		// this is a dirty trick to reuse this field to provide the password to the caller.
		HashedPassword: in.Password,
	}

	return u, nil
}

func verifyTOTPSecretForUser(ctx context.Context, c client.Client, userID, twoFactorSecret string) error {
	token, tokenErr := totp.GenerateCode(twoFactorSecret, time.Now().UTC())
	if tokenErr != nil {
		return fmt.Errorf("generating totp code: %w", tokenErr)
	}

	if _, err := c.VerifyTOTPSecret(ctx, &authsvc.VerifyTOTPSecretRequest{
		TOTPToken: token,
		UserID:    userID,
	}); err != nil {
		return fmt.Errorf("verifying totp code: %w", err)
	}

	return nil
}

func createClientForUser(ctx context.Context, scopes []string, user *identity.User) (client.Client, error) {
	token, err := fetchLoginTokenForUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("fetching token for user %s: %w", user.Username, err)
	}

	oauthedClient := buildAuthedGRPCClient(ctx, scopes, token)

	return oauthedClient, nil
}

func createUserAndClientForTest(t *testing.T) (*identity.User, client.Client) {
	t.Helper()

	ctx := t.Context()

	input := &identity.UserRegistrationInput{
		Birthday:              pointer.To(time.Now()),
		EmailAddress:          fmt.Sprintf("test+%d@whatever.com", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
		FirstName:             fmt.Sprintf("test_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
		AccountName:           fmt.Sprintf("test_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
		LastName:              fmt.Sprintf("test_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
		Password:              fmt.Sprintf("test_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
		Username:              fmt.Sprintf("test_%d", hashStringToNumber(t.Name()+time.Now().Format(time.RFC3339Nano))),
		AcceptedPrivacyPolicy: true,
		AcceptedTOS:           true,
	}

	user := createServiceUserForTest(t, true, input)
	oauthedClient := buildAuthedGRPCClient(ctx, []string{"account_admin"}, fetchLoginTokenForUserForTest(t, user))

	return user, oauthedClient
}

func fetchLoginTokenForUserForTest(t *testing.T, user *identity.User) string {
	t.Helper()
	ctx := t.Context()

	rv, err := fetchLoginTokenForUser(ctx, user)
	require.NoError(t, err)

	return rv
}

func generateTOTPCodeForUserForTest(t *testing.T, user *identity.User) string {
	t.Helper()

	code, err := totp.GenerateCode(strings.ToUpper(user.TwoFactorSecret), time.Now().UTC())
	require.NoError(t, err)

	return code
}

func generateTOTPCodeForUser(user *identity.User) (string, error) {
	code, err := totp.GenerateCode(strings.ToUpper(user.TwoFactorSecret), time.Now().UTC())
	if err != nil {
		return "", fmt.Errorf("generating totp code: %w", err)
	}

	return code, nil
}

func fetchLoginTokenForUser(ctx context.Context, user *identity.User) (string, error) {
	code, err := generateTOTPCodeForUser(user)
	if err != nil {
		return "", err
	}

	loginInput := &authsvc.UserLoginInput{
		Username:  user.Username,
		Password:  user.HashedPassword,
		TOTPToken: code,
	}

	// wretched hack that unfortunately works
	if user.Username == premadeAdminUser.Username {
		loginInput.Password = adminUserPassword
	}

	unauthedClient, err := buildUnauthenticatedGRPCClient()
	if err != nil {
		return "", fmt.Errorf("initializing client: %w", err)
	}

	tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
		Input: loginInput,
	})
	if err != nil {
		return "", fmt.Errorf("fetching login token: %w", err)
	}

	return tokenRes.Result.AccessToken, nil
}

//////// ChatGPT Zone

type compareOptions struct {
	// Ignore any field with these names at any depth (e.g., "LastUpdatedAt").
	IgnoreFieldNames map[string]struct{}
	// Only exported fields are considered (safe for cross-package types).
	ExportedOnly bool
}

// assertRoughEquality reports whether a and b are deeply equal after ignoring fields by name at any depth.
// Works across different struct types as long as exported field names/structure align.
func assertRoughEquality(t *testing.T, a, b any, ignoreFieldNames ...string) {
	t.Helper()

	opts := compareOptions{
		IgnoreFieldNames: toSet(ignoreFieldNames),
		ExportedOnly:     true,
	}
	ma := flattenComparable(a, opts)
	mb := flattenComparable(b, opts)
	diff := diffMaps(ma, mb)

	assert.True(t, len(diff) == 0, "objects should match except for LastUpdatedAt, diffs: %+v", diff)
}

func toSet(xs []string) map[string]struct{} {
	m := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		m[x] = struct{}{}
	}
	return m
}

// flattenComparable produces a deterministic, comparable map[path]string for any value.
// It skips fields listed in opts.IgnoreFieldNames (matched by the field name at any depth),
// only includes exported fields when opts.ExportedOnly is true, and handles cycles.
func flattenComparable(v any, opts compareOptions) map[string]string {
	out := make(map[string]string)
	visited := make(map[uintptr]struct{})
	var walk func(rv reflect.Value, path []string)

	shouldIgnoreField := func(fieldName string) bool {
		_, ok := opts.IgnoreFieldNames[fieldName]
		return ok
	}

	join := func(path []string) string {
		return strings.Join(path, ".")
	}

	// Handle time.Time specially for stable representation.
	writeTime := func(tv time.Time, path []string) {
		// Use RFC3339Nano for human-readable + stable, or UnixNano if you prefer strict numeric
		out[join(path)] = tv.UTC().Format(time.RFC3339Nano)
	}

	writeLeaf := func(rv reflect.Value, path []string) {
		// Convert leaf value to a stable string.
		switch rv.Kind() {
		case reflect.String:
			out[join(path)] = rv.String()
		case reflect.Bool:
			out[join(path)] = strconv.FormatBool(rv.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			out[join(path)] = strconv.FormatInt(rv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			out[join(path)] = strconv.FormatUint(rv.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			out[join(path)] = strconv.FormatFloat(rv.Float(), 'g', -1, rv.Type().Bits())
		case reflect.Complex64, reflect.Complex128:
			c := rv.Complex()
			out[join(path)] = fmt.Sprintf("(%g+%gi)", real(c), imag(c))
		default:
			// Fallback
			out[join(path)] = fmt.Sprintf("%v", rv.Interface())
		}
	}

	walk = func(rv reflect.Value, path []string) {
		if !rv.IsValid() {
			out[join(path)] = "nil"
			return
		}

		// Unwrap interfaces
		if rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				out[join(path)] = "nil"
				return
			}
			rv = rv.Elem()
		}

		// Follow pointers with cycle detection
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				out[join(path)] = "nil"
				return
			}
			ptr := rv.Pointer()
			if ptr != 0 {
				if _, seen := visited[ptr]; seen {
					// Prevent cycles
					out[join(path)] = "<cycle>"
					return
				}
				visited[ptr] = struct{}{}
			}
			walk(rv.Elem(), path)
			return
		}

		// time.Time special case
		if rv.Type() == reflect.TypeOf(time.Time{}) {
			writeTime(rv.Interface().(time.Time), path)
			return
		}

		switch rv.Kind() {
		case reflect.Struct:
			rt := rv.Type()
			for i := 0; i < rv.NumField(); i++ {
				sf := rt.Field(i)
				// Skip unexported fields if requested
				if opts.ExportedOnly && sf.PkgPath != "" {
					continue
				}
				if shouldIgnoreField(sf.Name) {
					continue
				}
				walk(rv.Field(i), append(path, sf.Name))
			}

		case reflect.Slice, reflect.Array:
			l := rv.Len()
			for i := 0; i < l; i++ {
				walk(rv.Index(i), append(path, fmt.Sprintf("[%d]", i)))
			}

		case reflect.Map:
			if rv.IsNil() {
				out[join(path)] = "nil"
				return
			}
			keys := rv.MapKeys()
			// Sort keys deterministically by their string form
			sort.Slice(keys, func(i, j int) bool {
				return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
			})
			for _, k := range keys {
				kv := rv.MapIndex(k)
				walk(kv, append(path, fmt.Sprintf("[%v]", k.Interface())))
			}

		default:
			// Basic leaf
			writeLeaf(rv, path)
		}
	}

	walk(reflect.ValueOf(v), nil)
	return out
}

func diffMaps(a, b map[string]string) map[string][2]string {
	diff := make(map[string][2]string)
	keys := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}
	for k := range keys {
		av, aok := a[k]
		bv, bok := b[k]
		switch {
		case aok && !bok:
			diff[k] = [2]string{av, "<absent>"}
		case !aok && bok:
			diff[k] = [2]string{"<absent>", bv}
		case aok && bok && av != bv:
			diff[k] = [2]string{av, bv}
		}
	}
	return diff
}
