package authorization

import (
	"github.com/mikespook/gorbac/v2"
)

const (
	// PublishArbitraryQueueMessagesPermission is a service admin permission.
	PublishArbitraryQueueMessagesPermission Permission = "admin.publish_queue_messages"
	// PublishArbitraryQueueMessagePermission is a service admin permission.
	PublishArbitraryQueueMessagePermission Permission = "queues.publish.message"
)

var (
	// AdminPermissions contains all admin-specific permissions.
	AdminPermissions = []gorbac.Permission{
		PublishArbitraryQueueMessagesPermission,
		PublishArbitraryQueueMessagePermission,
	}
)
