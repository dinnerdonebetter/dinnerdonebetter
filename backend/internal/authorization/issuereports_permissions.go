package authorization

const (
	// CreateIssueReportsPermission is an account admin permission.
	CreateIssueReportsPermission Permission = "create.issue_reports"
	// ReadIssueReportsPermission is an account admin permission.
	ReadIssueReportsPermission Permission = "read.issue_reports"
	// UpdateIssueReportsPermission is an account admin permission.
	UpdateIssueReportsPermission Permission = "update.issue_reports"
	// ArchiveIssueReportsPermission is an account admin permission.
	ArchiveIssueReportsPermission Permission = "archive.issue_reports"
)

var (
	// IssueReportsPermissions contains all issue report-related permissions.
	IssueReportsPermissions = []Permission{
		CreateIssueReportsPermission,
		ReadIssueReportsPermission,
		UpdateIssueReportsPermission,
		ArchiveIssueReportsPermission,
	}

	// IssueReportsAccountAdminPermissions contains issue report permissions for the account admin role.
	// Pass to RegisterAccountAdminPermissions in the domain registration module.
	IssueReportsAccountAdminPermissions = []Permission{
		CreateIssueReportsPermission,
		UpdateIssueReportsPermission,
		ArchiveIssueReportsPermission,
	}

	// IssueReportsAccountMemberPermissions contains issue report permissions for the account member role.
	// Pass to RegisterAccountMemberPermissions in the domain registration module.
	IssueReportsAccountMemberPermissions = []Permission{
		ReadIssueReportsPermission,
	}
)
