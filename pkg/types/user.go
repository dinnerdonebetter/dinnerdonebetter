package types

import (
	"context"
	"math"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	// UserDataType indicates an event is user membership-related.
	UserDataType dataType = "user"

	// UserSignedUpCustomerEventType indicates a user signed up.
	UserSignedUpCustomerEventType CustomerEventType = "user_signed_up"

	// GoodStandingUserAccountStatus indicates a User's household is in good standing.
	GoodStandingUserAccountStatus userAccountStatus = "good"
	// UnverifiedHouseholdStatus indicates a User's household requires two factor secret verification.
	UnverifiedHouseholdStatus userAccountStatus = "unverified"
	// BannedUserAccountStatus indicates a User's household is banned.
	BannedUserAccountStatus userAccountStatus = "banned"
	// TerminatedUserAccountStatus indicates a User's household is banned.
	TerminatedUserAccountStatus userAccountStatus = "terminated"

	validTOTPTokenLength = 6
)

var (
	totpTokenLengthRule = validation.Length(validTOTPTokenLength, validTOTPTokenLength)
)

type (
	userAccountStatus string

	// User represents a User.
	User struct {
		_                         struct{}
		PasswordLastChangedOn     *uint64           `json:"passwordLastChangedOn"`
		ArchivedOn                *uint64           `json:"archivedOn"`
		LastUpdatedOn             *uint64           `json:"lastUpdatedOn"`
		TwoFactorSecretVerifiedOn *uint64           `json:"-"`
		AvatarSrc                 *string           `json:"avatar"`
		BirthMonth                *uint8            `json:"birthMonth"`
		BirthDay                  *uint8            `json:"birthDay"`
		EmailAddress              string            `json:"emailAddress"`
		AccountStatusExplanation  string            `json:"accountStatusExplanation"`
		TwoFactorSecret           string            `json:"-"`
		HashedPassword            string            `json:"-"`
		ID                        string            `json:"id"`
		AccountStatus             userAccountStatus `json:"accountStatus"`
		Username                  string            `json:"username"`
		ServiceRoles              []string          `json:"serviceRoles"`
		CreatedOn                 uint64            `json:"createdOn"`
		RequiresPasswordChange    bool              `json:"requiresPasswordChange"`
	}

	// UserList represents a list of users.
	UserList struct {
		_ struct{}

		Users []*User `json:"data"`
		Pagination
	}

	// UserRegistrationInput represents the input required from users to register an account.
	UserRegistrationInput struct {
		_               struct{}
		BirthDay        *uint8 `json:"birthDay,omitempty"`
		BirthMonth      *uint8 `json:"birthMonth,omitempty"`
		Password        string `json:"password"`
		EmailAddress    string `json:"emailAddress"`
		InvitationToken string `json:"invitationToken,omitempty"`
		InvitationID    string `json:"invitationID,omitempty"`
		Username        string `json:"username"`
	}

	// UserDatabaseCreationInput is used by the User creation route to communicate with the data store.
	UserDatabaseCreationInput struct {
		_                    struct{}
		BirthMonth           *uint8  `json:"-"`
		BirthDay             *uint8  `json:"-"`
		ID                   string  `json:"-"`
		AvatarSrc            *string `json:"-"`
		HashedPassword       string  `json:"-"`
		TwoFactorSecret      string  `json:"-"`
		InvitationToken      string  `json:"-"`
		DestinationHousehold string  `json:"-"`
		Username             string  `json:"-"`
		EmailAddress         string  `json:"-"`
	}

	// UserCreationResponse is a response structure for Users that doesn't contain passwords fields, but does contain the two factor secret.
	UserCreationResponse struct {
		_               struct{}
		BirthMonth      *uint8            `json:"birthMonth"`
		BirthDay        *uint8            `json:"birthDay"`
		Username        string            `json:"username"`
		AvatarSrc       *string           `json:"avatar"`
		EmailAddress    string            `json:"emailAddress"`
		TwoFactorQRCode string            `json:"qrCode"`
		CreatedUserID   string            `json:"createdUserID"`
		HouseholdStatus userAccountStatus `json:"userAccountStatus"`
		TwoFactorSecret string            `json:"twoFactorSecret"`
		CreatedOn       uint64            `json:"createdOn"`
		IsAdmin         bool              `json:"isAdmin"`
	}

	// UserLoginInput represents the payload used to log in a User.
	UserLoginInput struct {
		_ struct{}

		Username  string `json:"username"`
		Password  string `json:"password"`
		TOTPToken string `json:"totpToken"`
	}

	// PasswordUpdateInput represents input a User would provide when updating their passwords.
	PasswordUpdateInput struct {
		_ struct{}

		NewPassword     string `json:"newPassword"`
		CurrentPassword string `json:"currentPassword"`
		TOTPToken       string `json:"totpToken"`
	}

	// TOTPSecretRefreshInput represents input a User would provide when updating their 2FA secret.
	TOTPSecretRefreshInput struct {
		_ struct{}

		CurrentPassword string `json:"currentPassword"`
		TOTPToken       string `json:"totpToken"`
	}

	// TOTPSecretVerificationInput represents input a User would provide when validating their 2FA secret.
	TOTPSecretVerificationInput struct {
		_ struct{}

		TOTPToken string `json:"totpToken"`
		UserID    string `json:"userID"`
	}

	// TOTPSecretRefreshResponse represents the response we provide to a User when updating their 2FA secret.
	TOTPSecretRefreshResponse struct {
		_ struct{}

		TwoFactorQRCode string `json:"qrCode"`
		TwoFactorSecret string `json:"twoFactorSecret"`
	}

	// AdminUserDataManager contains administrative User functions that we don't necessarily want to expose
	// to, say, the collection of handlers.
	AdminUserDataManager interface {
		UpdateUserAccountStatus(ctx context.Context, userID string, input *UserAccountStatusUpdateInput) error
	}

	// UserDataManager describes a structure which can manage users in permanent storage.
	UserDataManager interface {
		UserHasStatus(ctx context.Context, userID string, statuses ...string) (bool, error)
		GetUser(ctx context.Context, userID string) (*User, error)
		GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID string) (*User, error)
		GetUserIDByEmail(ctx context.Context, email string) (string, error)
		MarkUserTwoFactorSecretAsVerified(ctx context.Context, userID string) error
		GetUserByUsername(ctx context.Context, username string) (*User, error)
		GetAdminUserByUsername(ctx context.Context, username string) (*User, error)
		SearchForUsersByUsername(ctx context.Context, usernameQuery string) ([]*User, error)
		GetAllUsersCount(ctx context.Context) (uint64, error)
		GetUsers(ctx context.Context, filter *QueryFilter) (*UserList, error)
		CreateUser(ctx context.Context, input *UserDatabaseCreationInput) (*User, error)
		UpdateUser(ctx context.Context, updated *User) error
		UpdateUserPassword(ctx context.Context, userID, newHash string) error
		ArchiveUser(ctx context.Context, userID string) error
	}

	// UserDataService describes a structure capable of serving traffic related to users.
	UserDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		SelfHandler(res http.ResponseWriter, req *http.Request)
		UsernameSearchHandler(res http.ResponseWriter, req *http.Request)
		NewTOTPSecretHandler(res http.ResponseWriter, req *http.Request)
		TOTPSecretVerificationHandler(res http.ResponseWriter, req *http.Request)
		UpdatePasswordHandler(res http.ResponseWriter, req *http.Request)
		AvatarUploadHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update accepts a User as input and merges those values if they're set.
func (u *User) Update(input *User) {
	if input.Username != "" && input.Username != u.Username {
		u.Username = input.Username
	}

	if input.HashedPassword != "" && input.HashedPassword != u.HashedPassword {
		u.HashedPassword = input.HashedPassword
	}

	if input.TwoFactorSecret != "" && input.TwoFactorSecret != u.TwoFactorSecret {
		u.TwoFactorSecret = input.TwoFactorSecret
	}

	if input.BirthDay != nil && *input.BirthDay != 0 && input.BirthDay != u.BirthDay {
		u.BirthDay = input.BirthDay
	}

	if input.BirthMonth != nil && *input.BirthMonth != 0 && input.BirthMonth != u.BirthMonth {
		u.BirthMonth = input.BirthMonth
	}
}

// IsValidHouseholdStatus returns whether the provided string is a valid userAccountStatus.
func IsValidHouseholdStatus(s string) bool {
	switch s {
	case string(GoodStandingUserAccountStatus),
		string(UnverifiedHouseholdStatus),
		string(BannedUserAccountStatus),
		string(TerminatedUserAccountStatus):
		return true
	default:
		return false
	}
}

// IsBanned is a handy helper function.
func (u *User) IsBanned() bool {
	return u.AccountStatus == BannedUserAccountStatus
}

// ValidateWithContext ensures our provided UserRegistrationInput meets expectations.
func (i *UserRegistrationInput) ValidateWithContext(ctx context.Context, minUsernameLength, minPasswordLength uint8) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.EmailAddress, validation.Required, is.EmailFormat),
		validation.Field(&i.Username, validation.Required, validation.Length(int(minUsernameLength), math.MaxInt8)),
		validation.Field(&i.Password, validation.Required, validation.Length(int(minPasswordLength), math.MaxInt8)),
	)
}

var _ validation.ValidatableWithContext = (*TOTPSecretVerificationInput)(nil)

// ValidateWithContext ensures our provided TOTPSecretVerificationInput meets expectations.
func (i *TOTPSecretVerificationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.UserID, validation.Required),
		validation.Field(&i.TOTPToken, validation.Required, totpTokenLengthRule),
	)
}

// ValidateWithContext ensures our provided UserLoginInput meets expectations.
func (i *UserLoginInput) ValidateWithContext(ctx context.Context, minUsernameLength, minPasswordLength uint8) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.Username, validation.Required, validation.Length(int(minUsernameLength), math.MaxInt8)),
		validation.Field(&i.Password, validation.Required, validation.Length(int(minPasswordLength), math.MaxInt8)),
	)
}

// ValidateWithContext ensures our provided PasswordUpdateInput meets expectations.
func (i *PasswordUpdateInput) ValidateWithContext(ctx context.Context, minPasswordLength uint8) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.CurrentPassword, validation.Required, validation.Length(int(minPasswordLength), math.MaxInt8)),
		validation.Field(&i.NewPassword, validation.Required, validation.Length(int(minPasswordLength), math.MaxInt8)),
		validation.Field(&i.TOTPToken, validation.Required, totpTokenLengthRule),
	)
}

var _ validation.ValidatableWithContext = (*TOTPSecretRefreshInput)(nil)

// ValidateWithContext ensures our provided TOTPSecretRefreshInput meets expectations.
func (i *TOTPSecretRefreshInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.CurrentPassword, validation.Required),
		validation.Field(&i.TOTPToken, validation.Required, totpTokenLengthRule),
	)
}
