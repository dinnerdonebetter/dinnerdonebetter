package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/logging"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"

	"go.opentelemetry.io/otel/trace"

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

	client, err := httpclient.NewClient(parsedURI, trace.NewNoopTracerProvider(), httpclient.UsingCookie(cookie))
	if err != nil {
		return nil, nil, fmt.Errorf("initializing REST API client: %w", err)
	}

	return user, client, nil
}

var allDays = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
	time.Sunday,
}

var allMealNames = []types.MealName{
	types.BreakfastMealName,
	types.SecondBreakfastMealName,
	types.BrunchMealName,
	types.LunchMealName,
	types.SupperMealName,
	types.DinnerMealName,
}

func byDayAndMeal(l []*types.MealPlanOption, day time.Weekday, meal types.MealName) []*types.MealPlanOption {
	out := []*types.MealPlanOption{}

	for _, o := range l {
		if o.Day == day && o.MealName == meal {
			out = append(out, o)
		}
	}

	return out
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

			exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
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

func buildHandler(logger logging.Logger) func(ctx context.Context) error {
	return func(ctx context.Context) error {
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
				FromUser:             householdLeader.ID,
				Note:                 "prober testing",
				ToEmail:              u.EmailAddress,
				DestinationHousehold: relevantHouseholdID,
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

		exampleMealPlan := &types.MealPlan{
			Notes:          "prober testing",
			Status:         types.AwaitingVotesMealPlanStatus,
			StartsAt:       uint64(time.Now().Add(24 * time.Hour).Unix()),
			EndsAt:         uint64(time.Now().Add(72 * time.Hour).Unix()),
			VotingDeadline: uint64(time.Now().Add(votingDeadline).Unix()),
			Options: []*types.MealPlanOption{
				{
					MealID:   createdMeals[0].ID,
					Notes:    "option A",
					MealName: types.BreakfastMealName,
					Day:      time.Monday,
				},
				{
					MealID:   createdMeals[1].ID,
					Notes:    "option B",
					MealName: types.BreakfastMealName,
					Day:      time.Monday,
				},
				{
					MealID:   createdMeals[2].ID,
					Notes:    "option C",
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

		logger.Debug("fetched created meal plan")

		userAVotes := []*types.MealPlanOptionVote{
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
		}

		userBVotes := []*types.MealPlanOptionVote{
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
		}

		for _, vote := range userAVotes {
			logger.Debug("voting for user A")

			exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(vote)
			_, err = createdClients[0].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)

			if err != nil {
				return fmt.Errorf("voting for user A: %w", err)
			}
		}

		for _, vote := range userBVotes {
			logger.Debug("voting for user B")

			exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(vote)
			_, err = createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)

			if err != nil {
				return fmt.Errorf("voting for user B: %w", err)
			}
		}

		logger.Debug("getting voted upon meal plan")
		createdMealPlan, err = householdLeaderClient.GetMealPlan(ctx, createdMealPlan.ID)
		if err != nil {
			return fmt.Errorf("fetching voted upon meal plan: %w", err)
		}

		if types.AwaitingVotesMealPlanStatus != createdMealPlan.Status {
			return fmt.Errorf("unexpected meal plan status: %s", createdMealPlan.Status)
		}

		logger.Debug("waiting for worker to finalize meal plan")
		time.Sleep(creationDeadline)

		logger.Debug("getting hopefully finalized meal plan")
		createdMealPlan, err = householdLeaderClient.GetMealPlan(ctx, createdMealPlan.ID)
		if err != nil {
			return fmt.Errorf("fetching maybe finalized meal plan: %w", err)
		}

		if types.FinalizedMealPlanStatus != createdMealPlan.Status {
			return fmt.Errorf("unexpected final meal status: %s", createdMealPlan.Status)
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

		return nil
	}
}

func main() {
	logger := zerolog.NewZerologLogger()
	logger.SetLevel(logging.DebugLevel)
	logger.Info("starting prober")

	lambda.Start(buildHandler(logger))
}
