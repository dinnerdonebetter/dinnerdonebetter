package email

import (
	"errors"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/lib/email"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/matcornic/hermes/v2"
)

var (
	ErrUnverifiedEmailRecipient = errors.New("missing email address verification for user")
)

const (
	companyName = "Dinner Done Better"

	logoURL = "data:image/png;base64,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"
)

// BuildMealPlanCreatedEmail builds an email notifying a user that they've been invited to join a household.
func BuildMealPlanCreatedEmail(recipient *types.User, mealPlan *types.MealPlan, envCfg *email.EnvironmentConfig) (*email.OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, internalerrors.NilConfigError("email environment config")
	}

	if recipient.EmailAddressVerifiedAt == nil {
		return nil, ErrUnverifiedEmailRecipient
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: recipient.Username,
			Intros: []string{
				"A new meal plan has been created for your household!",
			},
			Actions: []hermes.Action{
				{
					Instructions: "You can rank each meal in the meal plan by clicking the button below",
					Button: hermes.Button{
						Text: "Submit your vote",
						Link: fmt.Sprintf("%s/meal_plans/%s", envCfg.BaseURL(), mealPlan.ID),
					},
				},
			},
		},
	}

	htmlContent, err := envCfg.BuildHermes().GenerateHTML(e)
	if err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &email.OutboundEmailMessage{
		ToAddress:   recipient.EmailAddress,
		FromAddress: envCfg.PasswordResetRedemptionEmailAddress(),
		FromName:    companyName,
		Subject:     "A new meal plan has been created!",
		HTMLContent: htmlContent,
	}

	return msg, nil
}
