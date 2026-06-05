package authorization

const (
	CreateProductsPermission  Permission = "create.products"
	ReadProductsPermission    Permission = "read.products"
	UpdateProductsPermission  Permission = "update.products"
	ArchiveProductsPermission Permission = "archive.products"

	CreateCheckoutSessionPermission Permission = "create.checkout_sessions"

	CreateSubscriptionsPermission  Permission = "create.subscriptions"
	ReadSubscriptionsPermission    Permission = "read.subscriptions"
	UpdateSubscriptionsPermission  Permission = "update.subscriptions"
	ArchiveSubscriptionsPermission Permission = "archive.subscriptions"
	CancelSubscriptionPermission   Permission = "cancel.subscriptions"

	ReadPurchasesPermission      Permission = "read.purchases"
	ReadPaymentHistoryPermission Permission = "read.payment_history"
)

var (
	// PaymentsPermissions contains all payments-related permissions.
	PaymentsPermissions = []Permission{
		CreateProductsPermission,
		ReadProductsPermission,
		UpdateProductsPermission,
		ArchiveProductsPermission,
		CreateCheckoutSessionPermission,
		CreateSubscriptionsPermission,
		ReadSubscriptionsPermission,
		UpdateSubscriptionsPermission,
		ArchiveSubscriptionsPermission,
		CancelSubscriptionPermission,
		ReadPurchasesPermission,
		ReadPaymentHistoryPermission,
	}

	// PaymentsServiceAdminPermissions contains payments permissions for the service admin role.
	// Pass to RegisterServiceAdminPermissions in the domain registration module.
	PaymentsServiceAdminPermissions = []Permission{
		CreateProductsPermission,
		ReadProductsPermission,
		UpdateProductsPermission,
		ArchiveProductsPermission,
		CreateSubscriptionsPermission,
		ReadSubscriptionsPermission,
		UpdateSubscriptionsPermission,
		ArchiveSubscriptionsPermission,
	}

	// PaymentsAccountAdminPermissions contains payments permissions for the account admin role.
	// Pass to RegisterAccountAdminPermissions in the domain registration module.
	PaymentsAccountAdminPermissions = []Permission{
		CreateCheckoutSessionPermission,
		CancelSubscriptionPermission,
		ReadPurchasesPermission,
		ReadPaymentHistoryPermission,
		ReadSubscriptionsPermission,
	}

	// PaymentsAccountMemberPermissions contains payments permissions for the account member role.
	// Pass to RegisterAccountMemberPermissions in the domain registration module.
	PaymentsAccountMemberPermissions = []Permission{
		CreateCheckoutSessionPermission,
		CancelSubscriptionPermission,
		ReadPurchasesPermission,
		ReadPaymentHistoryPermission,
		ReadSubscriptionsPermission,
	}
)
