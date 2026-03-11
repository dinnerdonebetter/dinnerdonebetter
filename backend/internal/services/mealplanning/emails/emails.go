package email

import (
	"errors"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/branding"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/email"
	"github.com/dinnerdonebetter/backend/internal/platform/internalerrors"

	"github.com/matcornic/hermes/v2"
)

var (
	ErrUnverifiedEmailRecipient = errors.New("missing email address verification for user")
)

// BuildMealPlanCreatedEmail builds an email notifying a user that they've been invited to join an account.
func BuildMealPlanCreatedEmail(recipient *identity.User, mealPlan *mealplanning.MealPlan, envCfg *email.EnvironmentConfig) (*email.OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, internalerrors.NilConfigError("email environment config")
	}

	if recipient.EmailAddressVerifiedAt == nil {
		return nil, ErrUnverifiedEmailRecipient
	}

	isElectionMealPlan := false
	for _, event := range mealPlan.Events {
		if len(event.Options) > 1 {
			isElectionMealPlan = true
		}
	}

	instructions := "You can see what's up for dinner by clicking the button below"
	if isElectionMealPlan {
		instructions = "You can rank each meal in the meal plan by clicking the button below"
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: recipient.FirstName,
			Intros: []string{
				"A new meal plan has been created for your account!",
			},
			Actions: []hermes.Action{
				{
					Instructions: instructions,
					Button: hermes.Button{
						Text: "Submit your vote",
						Link: fmt.Sprintf("%s/meal_plans/%s", envCfg.BaseURL(), mealPlan.ID),
					},
				},
			},
		},
	}

	htmlContent, err := envCfg.BuildHermes(branding.DefaultEmailBranding()).GenerateHTML(e)
	if err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &email.OutboundEmailMessage{
		ToAddress:   recipient.EmailAddress,
		FromAddress: envCfg.PasswordResetRedemptionEmailAddress(),
		FromName:    branding.CompanyName,
		Subject:     "A new meal plan has been created!",
		HTMLContent: htmlContent,
	}

	return msg, nil
}
