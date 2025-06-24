package identity

type Repository interface {
	AuditLogEntryDataManager
	AccountDataManager
	AccountInvitationDataManager
	PasswordResetTokenDataManager
	UserDataManager
}
