package users

import "net/http"

type (
	// UserService describes a structure capable of serving traffic related to users.
	UserService interface {
		ListUsersHandler(http.ResponseWriter, *http.Request)
		CreateUserHandler(http.ResponseWriter, *http.Request)
		ReadUserHandler(http.ResponseWriter, *http.Request)
		SelfHandler(http.ResponseWriter, *http.Request)
		UserPermissionsHandler(http.ResponseWriter, *http.Request)
		UsernameSearchHandler(http.ResponseWriter, *http.Request)
		NewTOTPSecretHandler(http.ResponseWriter, *http.Request)
		TOTPSecretVerificationHandler(http.ResponseWriter, *http.Request)
		UpdatePasswordHandler(http.ResponseWriter, *http.Request)
		UpdateUserEmailAddressHandler(http.ResponseWriter, *http.Request)
		UpdateUserUsernameHandler(http.ResponseWriter, *http.Request)
		UpdateUserDetailsHandler(http.ResponseWriter, *http.Request)
		AvatarUploadHandler(http.ResponseWriter, *http.Request)
		ArchiveUserHandler(http.ResponseWriter, *http.Request)
		CreatePasswordResetTokenHandler(http.ResponseWriter, *http.Request)
		PasswordResetTokenRedemptionHandler(http.ResponseWriter, *http.Request)
		RequestUsernameReminderHandler(http.ResponseWriter, *http.Request)
		VerifyUserEmailAddressHandler(http.ResponseWriter, *http.Request)
		RequestEmailVerificationEmailHandler(http.ResponseWriter, *http.Request)
	}
)
