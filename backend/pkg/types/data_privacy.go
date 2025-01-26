package types

import (
	"context"
	"net/http"
)

type (
	// DataDeletionResponse is returned when a user requests their data be deleted.
	DataDeletionResponse struct {
		Successful bool `json:"Successful"`
	}

	// UserDataAggregationRequest represents a message queue event meant to aggregate data for a user.
	UserDataAggregationRequest struct {
		_ struct{} `json:"-"`

		RequestID string `json:"id"`
		ReportID  string `json:"reportID"`
		UserID    string `json:"userID"`
	}

	// UserDataCollectionResponse represents the response to a UserDataAggregationRequest.
	UserDataCollectionResponse struct {
		_ struct{} `json:"-"`

		ReportID string `json:"reportID"`
	}

	UserDataCollection struct {
		_ struct{} `json:"-"`

		Core     CoreUserDataCollection   `json:"core"`
		Eating   EatingUserDataCollection `json:"eating"`
		ReportID string                   `json:"reportID"`
		User     User                     `json:"user"`
	}

	CoreUserDataCollection struct {
		_ struct{} `json:"-"`

		Webhooks                         map[string][]Webhook                     `json:"webhooks"`
		ServiceSettingConfigurations     map[string][]ServiceSettingConfiguration `json:"serviceSettingConfigurations"`
		AuditLogEntries                  map[string][]AuditLogEntry               `json:"auditLogEntries"`
		ReceivedInvites                  []HouseholdInvitation                    `json:"receivedInvites"`
		SentInvites                      []HouseholdInvitation                    `json:"sentInvites"`
		UserServiceSettingConfigurations []ServiceSettingConfiguration            `json:"userServiceSettingConfigurations"`
		UserAuditLogEntries              []AuditLogEntry                          `json:"userAuditLogEntries"`
		Households                       []Household                              `json:"households"`
	}

	EatingUserDataCollection struct {
		_ struct{} `json:"-"`

		HouseholdInstrumentOwnerships map[string][]HouseholdInstrumentOwnership `json:"householdInstrumentOwnerships"`
		MealPlans                     map[string][]MealPlan                     `json:"mealPlans"`
		ReportID                      string                                    `json:"reportID"`
		RecipeRatings                 []RecipeRating                            `json:"recipeRatings"`
		Recipes                       []Recipe                                  `json:"recipes"`
		Meals                         []Meal                                    `json:"meals"`
		UserIngredientPreferences     []UserIngredientPreference                `json:"userIngredientPreferences"`
	}

	// DataPrivacyService describes a structure capable of serving CCPA/GRPC-related requests.
	DataPrivacyService interface {
		DataDeletionHandler(http.ResponseWriter, *http.Request)
		UserDataAggregationRequestHandler(http.ResponseWriter, *http.Request)
		ReadUserDataAggregationReportHandler(http.ResponseWriter, *http.Request)
	}

	// DataPrivacyDataManager contains data privacy management functions.
	DataPrivacyDataManager interface {
		DeleteUser(ctx context.Context, userID string) error
		AggregateUserData(ctx context.Context, userID string) (*UserDataCollection, error)
	}
)
