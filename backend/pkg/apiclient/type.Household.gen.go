// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	Household struct {
		CreatedAt                string                            `json:"createdAt"`
		ZipCode                  string                            `json:"zipCode"`
		AddressLine1             string                            `json:"addressLine1"`
		BelongsToUser            string                            `json:"belongsToUser"`
		BillingStatus            string                            `json:"billingStatus"`
		City                     string                            `json:"city"`
		ContactPhone             string                            `json:"contactPhone"`
		Country                  string                            `json:"country"`
		ArchivedAt               string                            `json:"archivedAt"`
		AddressLine2             string                            `json:"addressLine2"`
		ID                       string                            `json:"id"`
		LastUpdatedAt            string                            `json:"lastUpdatedAt"`
		SubscriptionPlanID       string                            `json:"subscriptionPlanID"`
		State                    string                            `json:"state"`
		Name                     string                            `json:"name"`
		PaymentProcessorCustomer string                            `json:"paymentProcessorCustomer"`
		Members                  []HouseholdUserMembershipWithUser `json:"members"`
		Longitude                float64                           `json:"longitude"`
		Latitude                 float64                           `json:"latitude"`
	}
)
