package mealplanprober

import (
	"context"
	"fmt"
	"net/url"
	"time"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

const (
	stagingAddress = "https://api.prixfixe.dev"

	votingDeadline   = 45 * time.Second
	creationDeadline = 90 * time.Second
)

func getClientForUser(ctx context.Context, logger logging.Logger) (*types.User, *httpclient.Client, error) {
	example := fakes.BuildFakeUserRegistrationInput()
	input := &types.UserRegistrationInput{
		Username:     example.Username,
		Password:     example.Password,
		EmailAddress: example.EmailAddress,
	}

	parsedURI, err := url.Parse(stagingAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing provided URI: %w", err)
	}

	logger.Debug("creating user")

	user, err := testutils.CreateServiceUser(ctx, stagingAddress, input)
	if err != nil {
		return nil, nil, err
	}

	logger = logger.WithValue("username", user.Username)
	logger.Debug("getting login cookie")

	cookie, err := testutils.GetLoginCookie(ctx, stagingAddress, user)
	if err != nil {
		return nil, nil, fmt.Errorf("getting login cookie: %w", err)
	}

	logger.Debug("initializing client")

	client, err := httpclient.NewClient(parsedURI, tracing.NewNoopTracerProvider(), httpclient.UsingCookie(cookie))
	if err != nil {
		return nil, nil, fmt.Errorf("initializing REST API client: %w", err)
	}

	return user, client, nil
}

func stringPointer(s string) *string {
	return &s
}

func createRecipeForTest(ctx context.Context, client *httpclient.Client) ([]*types.ValidIngredient, *types.ValidPreparation, *types.Recipe, error) {
	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)

	createdValidPreparation, err := client.CreateValidPreparation(ctx, exampleValidPreparationInput)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating valid preparation: %w", err)
	}

	exampleValidInstrument := fakes.BuildFakeValidInstrument()
	exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
	createdValidInstrument, err := client.CreateValidInstrument(ctx, exampleValidInstrumentInput)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating valid instrument: %w", err)
	}

	exampleRecipe := fakes.BuildFakeRecipe()

	createdValidIngredients := []*types.ValidIngredient{}
	for i, recipeStep := range exampleRecipe.Steps {
		for j := range recipeStep.Ingredients {
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, validIngredientCreationErr := client.CreateValidIngredient(ctx, exampleValidIngredientInput)
			if validIngredientCreationErr != nil {
				return nil, nil, nil, fmt.Errorf("creating valid ingredient: %w", validIngredientCreationErr)
			}

			createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

			exampleRecipe.Steps[i].Ingredients[j].Ingredient.ID = createdValidIngredient.ID
		}

		for j := range recipeStep.Instruments {
			recipeStep.Instruments[j].Instrument = createdValidInstrument
		}
	}

	exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
	}

	createdRecipe, err := client.CreateRecipe(ctx, exampleRecipeInput)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating recipe: %w", err)
	}

	return createdValidIngredients, createdValidPreparation, createdRecipe, nil
}

func createMealForTest(ctx context.Context, client *httpclient.Client) (*types.Meal, error) {
	createdRecipeIDs := []string{}
	for i := 0; i < 3; i++ {
		_, _, recipe, err := createRecipeForTest(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("creating recipe #%d: %w", i, err)
		}
		createdRecipeIDs = append(createdRecipeIDs, recipe.ID)
	}

	exampleMeal := fakes.BuildFakeMeal()
	exampleMealInput := fakes.BuildFakeMealCreationRequestInputFromMeal(exampleMeal)
	exampleMealInput.Recipes = createdRecipeIDs

	createdMeal, err := client.CreateMeal(ctx, exampleMealInput)
	if err != nil {
		return nil, fmt.Errorf("creating meal: %w", err)
	}

	return createdMeal, nil
}

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ProbeMealPlanning is our cloud function entrypoint.
func ProbeMealPlanning(ctx context.Context, m PubSubMessage) error {
	/*
		logger := zerolog.NewZerologLogger()

		logger.Debug("creating household leader")
		householdLeader, householdLeaderClient, err := getClientForUser(ctx, logger)
		if err != nil {
			return fmt.Errorf("creating household leader API client")
		}

		logger.Debug("household leader created,determining household ID")
		currentStatus, err := householdLeaderClient.UserStatus(ctx)
		if err != nil {
			return fmt.Errorf("checking household leader user status: %w", err)
		}

		relevantHouseholdID := currentStatus.ActiveHousehold
		createdClients := []*httpclient.Client{}

		for i := 0; i < 2; i++ {
			logger.WithValue("i", i).Debug("creating user to invite")
			u, c, userCreationErr := getClientForUser(ctx, logger)
			if userCreationErr != nil {
				return fmt.Errorf("error creating household member #%d: %w", i, userCreationErr)
			}

			logger.WithValue("i", i).Debug("inviting user")
			invitation, invitationCreationErr := householdLeaderClient.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
				FromUser:               householdLeader.ID,
				Note:                   "prober testing",
				ToEmail:                u.EmailAddress,
				DestinationHouseholdID: relevantHouseholdID,
			})
			if invitationCreationErr != nil {
				return fmt.Errorf("inviting user #%d: %w", i, invitationCreationErr)
			}

			logger.WithValue("i", i).Debug("checking for sent invitation")
			sentInvitations, invitationSendErr := householdLeaderClient.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			if invitationSendErr != nil {
				return fmt.Errorf("checking for sent invitations for user #%d: %w", i, invitationSendErr)
			}

			if len(sentInvitations.HouseholdInvitations) == 0 {
				return fmt.Errorf("no invitations sent to user #%d", i)
			}

			logger.WithValue("i", i).Debug("checking for received invitation")
			invitations, getInvitationsErr := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			if getInvitationsErr != nil {
				return fmt.Errorf("checking for received invitations for user #%d: %w", i, getInvitationsErr)
			}

			if len(invitations.HouseholdInvitations) == 0 {
				return fmt.Errorf("user #%d received no invitations", i)
			}

			logger.WithValue("i", i).Debug("accepting household invitation")
			if err = c.AcceptHouseholdInvitation(ctx, relevantHouseholdID, invitation.ID, "prober testing"); err != nil {
				return fmt.Errorf("accepting household invitation for user #%d: %w", i, err)
			}

			logger.WithValue("i", i).Debug("switching active household")
			if err = c.SwitchActiveHousehold(ctx, relevantHouseholdID); err != nil {
				return fmt.Errorf("switching household for user #%d: %w", i, err)
			}
			createdClients = append(createdClients, c)
		}

		logger.Debug("creating recipes for meal plan")
		createdMeals := []*types.Meal{}
		for i := 0; i < 3; i++ {
			createdMeal, mealCreationErr := createMealForTest(ctx, householdLeaderClient)
			if mealCreationErr != nil {
				return fmt.Errorf("error creating meal #%d: %w", i, mealCreationErr)
			}
			createdMeals = append(createdMeals, createdMeal)
		}

		now := time.Now()
		exampleMealPlan := &types.MealPlan{
			StatusExplanation:          "prober testing",
			CreationExplanation:         types.AwaitingVotesMealPlanStatus,
			StartsAt:       now.Add(24 * time.Hour),
			EndsAt:         now.Add(72 * time.Hour),
			VotingDeadline: now.Add(votingDeadline),
			Options: []*types.MealPlanOption{
				{
					Meal:     types.Meal{ID: createdMeals[0].ID},
					StatusExplanation:    "option A",
					MealName: types.BreakfastMealName,
					Day:      time.Monday,
				},
				{
					Meal:     types.Meal{ID: createdMeals[1].ID},
					StatusExplanation:    "option B",
					MealName: types.BreakfastMealName,
					Day:      time.Monday,
				},
				{
					Meal:     types.Meal{ID: createdMeals[2].ID},
					StatusExplanation:    "option C",
					MealName: types.BreakfastMealName,
					Day:      time.Monday,
				},
			},
		}

		exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
		createdMealPlan, err := householdLeaderClient.CreateMealPlan(ctx, exampleMealPlanInput)
		if err != nil {
			return fmt.Errorf("creating meal plan: %w", err)
		}

		logger.Debug("created meal plan")

		createdMealPlan, err = householdLeaderClient.GetMealPlan(ctx, createdMealPlan.ID)
		if err != nil {
			return fmt.Errorf("fetching meal plan: %w", err)
		}

		logger.Debug("fetched created meal plan, voting for user A")

		exampleMealPlanOptionVoteInputA := &types.MealPlanOptionVoteCreationRequestInput{
			Votes: []*types.MealPlanOptionVoteCreationInput{
				{
					BelongsToMealPlanOption: createdMealPlan.Options[0].ID,
					Rank:                    0,
				},
				{
					BelongsToMealPlanOption: createdMealPlan.Options[1].ID,
					Rank:                    2,
				},
				{
					BelongsToMealPlanOption: createdMealPlan.Options[2].ID,
					Rank:                    1,
				},
			},
		}

		_, err = createdClients[0].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInputA)
		if err != nil {
			return fmt.Errorf("voting for user A: %w", err)
		}

		exampleMealPlanOptionVoteInputB := &types.MealPlanOptionVoteCreationRequestInput{
			Votes: []*types.MealPlanOptionVoteCreationInput{
				{
					BelongsToMealPlanOption: createdMealPlan.Options[0].ID,
					Rank:                    0,
				},
				{
					BelongsToMealPlanOption: createdMealPlan.Options[1].ID,
					Rank:                    1,
				},
				{
					BelongsToMealPlanOption: createdMealPlan.Options[2].ID,
					Rank:                    2,
				},
			},
		}

		logger.Debug("voting for user B")

		_, err = createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInputB)
		if err != nil {
			return fmt.Errorf("voting for user B: %w", err)
		}

		logger.Debug("getting voted upon meal plan")
		createdMealPlan, err = householdLeaderClient.GetMealPlan(ctx, createdMealPlan.ID)
		if err != nil {
			return fmt.Errorf("fetching voted upon meal plan: %w", err)
		}

		if types.AwaitingVotesMealPlanStatus != createdMealPlan.CreationExplanation {
			return fmt.Errorf("unexpected meal plan status: %s", createdMealPlan.CreationExplanation)
		}

		logger.Debug("waiting for worker to finalize meal plan")
		time.Sleep(creationDeadline)

		logger.Debug("getting hopefully finalized meal plan")
		createdMealPlan, err = householdLeaderClient.GetMealPlan(ctx, createdMealPlan.ID)
		if err != nil {
			return fmt.Errorf("fetching maybe finalized meal plan: %w", err)
		}

		if types.FinalizedMealPlanStatus != createdMealPlan.CreationExplanation {
			return fmt.Errorf("unexpected final meal status: %s", createdMealPlan.CreationExplanation)
		}

		for _, day := range allDays {
			for _, mealName := range allMealNames {
				options := byDayAndMeal(createdMealPlan.Options, day, mealName)
				if len(options) > 0 {
					selectionMade := false
					for _, opt := range options {
						if opt.Chosen {
							selectionMade = true
							break
						}
					}

					if !selectionMade {
						return fmt.Errorf("selection wasn't made for meal plan %s", createdMealPlan.ID)
					}
				}
			}
		}

		logger.Info("done")
	*/
	return nil
}
