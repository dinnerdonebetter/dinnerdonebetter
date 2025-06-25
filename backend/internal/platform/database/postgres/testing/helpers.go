package testing

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/implementations/identity/generated"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/matryer/try.v1"
)

var RunContainerTests = strings.ToLower(os.Getenv("RUN_CONTAINER_TESTS")) != "false" // on by default

func HashStringToNumber(t *testing.T, s string) uint64 {
	t.Helper()
	h := fnv.New64a()

	_, err := h.Write([]byte(s))
	require.NoError(t, err)

	return h.Sum64()
}

func reverseString(input string) string {
	runes := []rune(input)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func splitReverseConcat(input string) string {
	length := len(input)
	halfLength := length / 2

	firstHalf := input[:halfLength]
	secondHalf := input[halfLength:]

	reversedFirstHalf := reverseString(firstHalf)
	reversedSecondHalf := reverseString(secondHalf)

	return reversedSecondHalf + reversedFirstHalf
}

const (
	defaultPostgresImage = "postgres:17"
)

func BuildDatabaseClientForTest(t *testing.T) (*postgres.PostgresContainer, *sql.DB) {
	t.Helper()

	dbUsername := fmt.Sprintf("%d", HashStringToNumber(t, t.Name()))
	testcontainers.Logger = log.New(io.Discard, "", log.LstdFlags)

	ctx := t.Context()

	var container *postgres.PostgresContainer
	err := try.Do(func(attempt int) (bool, error) {
		var containerErr error
		container, containerErr = postgres.Run(
			ctx,
			defaultPostgresImage,
			postgres.WithDatabase(splitReverseConcat(dbUsername)),
			postgres.WithUsername(dbUsername),
			postgres.WithPassword(reverseString(dbUsername)),
			testcontainers.WithWaitStrategyAndDeadline(2*time.Minute, wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
		)

		return attempt < 5, containerErr
	})
	require.NoError(t, err)
	require.NotNil(t, container)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbConfig := &databasecfg.Config{
		RunMigrations:            true,
		OAuth2TokenEncryptionKey: "blahblahblahblahblahblahblahblah",
	}
	require.NoError(t, dbConfig.LoadConnectionDetailsFromURL(connStr))

	db, err := dbConfig.ConnectToDatabase()
	require.NoError(t, err)

	return container, db
}

func CreateUserForTest(t *testing.T, ctx context.Context, exampleUser *identity.User, db *sql.DB) *identity.User {
	t.Helper()

	// create
	if exampleUser == nil {
		exampleUser = fakes.BuildFakeUser()
	}
	exampleUser.TwoFactorSecretVerifiedAt = nil

	dbc := generated.New()

	err := dbc.CreateUser(ctx, db, &generated.CreateUserParams{
		ID:                            exampleUser.ID,
		Username:                      exampleUser.Username,
		AvatarSrc:                     database.NullStringFromStringPointer(exampleUser.AvatarSrc),
		EmailAddress:                  exampleUser.EmailAddress,
		HashedPassword:                exampleUser.HashedPassword,
		RequiresPasswordChange:        exampleUser.RequiresPasswordChange,
		TwoFactorSecret:               exampleUser.TwoFactorSecret,
		TwoFactorSecretVerifiedAt:     database.NullTimeFromTimePointer(exampleUser.TwoFactorSecretVerifiedAt),
		ServiceRole:                   exampleUser.ServiceRole,
		UserAccountStatus:             exampleUser.AccountStatus,
		UserAccountStatusExplanation:  exampleUser.AccountStatusExplanation,
		Birthday:                      database.NullTimeFromTimePointer(exampleUser.Birthday),
		EmailAddressVerificationToken: database.NullStringFromString("token"),
		FirstName:                     exampleUser.FirstName,
		LastName:                      exampleUser.LastName,
	})
	require.NoError(t, err)

	dbCreated, err := dbc.GetUserByID(ctx, db, exampleUser.ID)
	require.NoError(t, err)

	created := &identity.User{
		CreatedAt:                  dbCreated.CreatedAt,
		PasswordLastChangedAt:      database.TimePointerFromNullTime(dbCreated.PasswordLastChangedAt),
		LastUpdatedAt:              database.TimePointerFromNullTime(dbCreated.LastUpdatedAt),
		LastAcceptedTermsOfService: database.TimePointerFromNullTime(dbCreated.LastAcceptedTermsOfService),
		LastAcceptedPrivacyPolicy:  database.TimePointerFromNullTime(dbCreated.LastAcceptedPrivacyPolicy),
		TwoFactorSecretVerifiedAt:  database.TimePointerFromNullTime(dbCreated.TwoFactorSecretVerifiedAt),
		AvatarSrc:                  database.StringPointerFromNullString(dbCreated.AvatarSrc),
		Birthday:                   database.TimePointerFromNullTime(dbCreated.Birthday),
		ArchivedAt:                 database.TimePointerFromNullTime(dbCreated.ArchivedAt),
		AccountStatusExplanation:   dbCreated.UserAccountStatusExplanation,
		TwoFactorSecret:            dbCreated.TwoFactorSecret,
		HashedPassword:             dbCreated.HashedPassword,
		ID:                         dbCreated.ID,
		AccountStatus:              dbCreated.UserAccountStatus,
		Username:                   dbCreated.Username,
		FirstName:                  dbCreated.FirstName,
		LastName:                   dbCreated.LastName,
		EmailAddress:               dbCreated.EmailAddress,
		EmailAddressVerifiedAt:     database.TimePointerFromNullTime(dbCreated.EmailAddressVerifiedAt),
		ServiceRole:                dbCreated.ServiceRole,
		RequiresPasswordChange:     dbCreated.RequiresPasswordChange,
	}
	exampleUser.CreatedAt = created.CreatedAt
	exampleUser.Birthday = created.Birthday
	exampleUser.TwoFactorSecretVerifiedAt = created.TwoFactorSecretVerifiedAt
	assert.Equal(t, exampleUser, created)

	return created
}

func CreateAccountForTest(t *testing.T, ctx context.Context, exampleAccount *identity.Account, userID string, db *sql.DB) *identity.Account {
	t.Helper()

	// create
	if exampleAccount == nil {
		exampleAccount = fakes.BuildFakeAccount()
		exampleAccount.BelongsToUser = userID
	}
	exampleAccount.PaymentProcessorCustomerID = ""
	exampleAccount.Members = nil

	dbc := generated.New()

	require.NoError(t, dbc.CreateAccount(ctx, db, &generated.CreateAccountParams{
		ID:                exampleAccount.ID,
		Name:              exampleAccount.Name,
		BillingStatus:     exampleAccount.BillingStatus,
		ContactPhone:      exampleAccount.ContactPhone,
		BelongsToUser:     exampleAccount.BelongsToUser,
		AddressLine1:      exampleAccount.AddressLine1,
		AddressLine2:      exampleAccount.AddressLine2,
		City:              exampleAccount.City,
		State:             exampleAccount.State,
		ZipCode:           exampleAccount.ZipCode,
		Country:           exampleAccount.Country,
		Latitude:          database.NullStringFromFloat64Pointer(exampleAccount.Latitude),
		Longitude:         database.NullStringFromFloat64Pointer(exampleAccount.Longitude),
		WebhookHmacSecret: exampleAccount.WebhookEncryptionKey,
	}))

	dbCreated, err := dbc.GetAccountsForUser(ctx, db, &generated.GetAccountsForUserParams{
		BelongsToUser: userID,
	})
	require.NoError(t, err)
	require.Len(t, dbCreated, 1)

	created := &identity.Account{
		CreatedAt:                  dbCreated[0].CreatedAt,
		SubscriptionPlanID:         database.StringPointerFromNullString(dbCreated[0].SubscriptionPlanID),
		LastUpdatedAt:              database.TimePointerFromNullTime(dbCreated[0].LastUpdatedAt),
		ArchivedAt:                 database.TimePointerFromNullTime(dbCreated[0].ArchivedAt),
		Longitude:                  database.Float64PointerFromNullString(dbCreated[0].Longitude),
		Latitude:                   database.Float64PointerFromNullString(dbCreated[0].Latitude),
		State:                      dbCreated[0].State,
		ContactPhone:               dbCreated[0].ContactPhone,
		City:                       dbCreated[0].City,
		AddressLine1:               dbCreated[0].AddressLine1,
		ZipCode:                    dbCreated[0].ZipCode,
		Country:                    dbCreated[0].Country,
		BillingStatus:              dbCreated[0].BillingStatus,
		AddressLine2:               dbCreated[0].AddressLine2,
		PaymentProcessorCustomerID: dbCreated[0].PaymentProcessorCustomerID,
		BelongsToUser:              dbCreated[0].BelongsToUser,
		ID:                         dbCreated[0].ID,
		Name:                       dbCreated[0].Name,
		WebhookEncryptionKey:       dbCreated[0].WebhookHmacSecret,
	}

	exampleAccount.CreatedAt = created.CreatedAt
	exampleAccount.WebhookEncryptionKey = created.WebhookEncryptionKey
	assert.Equal(t, exampleAccount, created)

	return created
}
