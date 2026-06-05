package authorization

const (
	// CreateCommentsPermission is a permission.
	CreateCommentsPermission Permission = "create.comments"
	// ReadCommentsPermission is a permission.
	ReadCommentsPermission Permission = "read.comments"
	// UpdateCommentsPermission is a permission.
	UpdateCommentsPermission Permission = "update.comments"
	// ArchiveCommentsPermission is a permission.
	ArchiveCommentsPermission Permission = "archive.comments"
)

var (
	// CommentsPermissions contains all comment-related permissions.
	CommentsPermissions = []Permission{
		CreateCommentsPermission,
		ReadCommentsPermission,
		UpdateCommentsPermission,
		ArchiveCommentsPermission,
	}

	// CommentsAccountMemberPermissions contains comment permissions for the account member role.
	// Pass to RegisterAccountMemberPermissions in the domain registration module.
	CommentsAccountMemberPermissions = []Permission{
		CreateCommentsPermission,
		ReadCommentsPermission,
		UpdateCommentsPermission,
		ArchiveCommentsPermission,
	}
)
