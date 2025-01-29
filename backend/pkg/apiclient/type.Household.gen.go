// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
Household struct {
   AddressLine1 string `json:"addressLine1"`
 AddressLine2 string `json:"addressLine2"`
 ArchivedAt string `json:"archivedAt"`
 BelongsToUser string `json:"belongsToUser"`
 BillingStatus string `json:"billingStatus"`
 City string `json:"city"`
 ContactPhone string `json:"contactPhone"`
 Country string `json:"country"`
 CreatedAt string `json:"createdAt"`
 ID string `json:"id"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Latitude float64 `json:"latitude"`
 Longitude float64 `json:"longitude"`
 Members []HouseholdUserMembershipWithUser `json:"members"`
 Name string `json:"name"`
 PaymentProcessorCustomer string `json:"paymentProcessorCustomer"`
 State string `json:"state"`
 SubscriptionPlanID string `json:"subscriptionPlanID"`
 ZipCode string `json:"zipCode"`

}
)
