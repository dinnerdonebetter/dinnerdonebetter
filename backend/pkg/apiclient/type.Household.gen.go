// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	Household struct {
		CreatedAt                string                            `json:"createdAt"`
		LastUpdatedAt            string                            `json:"lastUpdatedAt"`
		ArchivedAt               string                            `json:"archivedAt"`
		ID                       string                            `json:"id"`
		BillingStatus            string                            `json:"billingStatus"`
		City                     string                            `json:"city"`
		ContactPhone             string                            `json:"contactPhone"`
		Country                  string                            `json:"country"`
		AddressLine2             string                            `json:"addressLine2"`
		AddressLine1             string                            `json:"addressLine1"`
		BelongsToUser            string                            `json:"belongsToUser"`
		ZipCode                  string                            `json:"zipCode"`
		SubscriptionPlanID       string                            `json:"subscriptionPlanID"`
		State                    string                            `json:"state"`
		Name                     string                            `json:"name"`
		PaymentProcessorCustomer string                            `json:"paymentProcessorCustomer"`
		Members                  []HouseholdUserMembershipWithUser `json:"members"`
		Longitude                float64                           `json:"longitude"`
		Latitude                 float64                           `json:"latitude"`
	}
)
