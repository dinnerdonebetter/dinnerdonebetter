package auth

type Repository interface {
	AuditLogEntryDataManager
	AccountDataManager
	AccountInvitationDataManager
	PasswordResetTokenDataManager
	UserDataManager
}
