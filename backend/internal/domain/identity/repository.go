package identity

type Repository interface {
	AccountDataManager
	AccountInvitationDataManager
	PasswordResetTokenDataManager
	UserDataManager
	AccountUserMembershipDataManager
}
