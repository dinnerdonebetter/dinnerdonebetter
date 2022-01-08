package main

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/logging"

	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

const (
	stagingAddress = "https://api.prixfixe.dev"

	creationDeadline = 1 * time.Minute
)

func getClientForUser(ctx context.Context) (*types.User, *httpclient.Client) {
	example := fakes.BuildFakeUserRegistrationInput()
	input := &types.UserRegistrationInput{
		Username:     example.Username,
		Password:     example.Password,
		EmailAddress: example.EmailAddress,
	}

	parsedURI, err := url.Parse(stagingAddress)
	if err != nil {
		panic(err)
	}

	user, err := testutils.CreateServiceUser(ctx, stagingAddress, input)
	if err != nil {
		panic(err)
	}

	cookie, err := testutils.GetLoginCookie(ctx, stagingAddress, user)
	if err != nil {
		panic(err)
	}

	client, err := httpclient.NewClient(parsedURI, trace.NewNoopTracerProvider(), httpclient.UsingCookie(cookie))
	if err != nil {
		panic(err)
	}

	return user, client
}

func mustnt(err error, doing string) {
	if err != nil {
		log.Panicf("error %s: %v", doing, err)
	}
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

func createRecipeForTest(ctx context.Context, logger logging.Logger, client *httpclient.Client) ([]*types.ValidIngredient, *types.ValidPreparation, *types.Recipe) {
	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
	createdValidPreparation, err := client.CreateValidPreparation(ctx, exampleValidPreparationInput)
	mustnt(err, "creating valid preparation")

	exampleRecipe := fakes.BuildFakeRecipe()

	createdValidIngredients := []*types.ValidIngredient{}
	for i, recipeStep := range exampleRecipe.Steps {
		for j := range recipeStep.Ingredients {
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := client.CreateValidIngredient(ctx, exampleValidIngredientInput)
			mustnt(err, "creating valid ingredient")

			createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

			exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
		}
	}

	exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
	}

	createdRecipe, err := client.CreateRecipe(ctx, exampleRecipeInput)
	mustnt(err, "creating recipe")

	return createdValidIngredients, createdValidPreparation, createdRecipe
}

func createMealForTest(ctx context.Context, logger logging.Logger, client *httpclient.Client) *types.Meal {
	createdRecipes := []*types.Recipe{}
	createdRecipeIDs := []string{}
	for i := 0; i < 3; i++ {
		_, _, recipe := createRecipeForTest(ctx, logger, client)
		createdRecipes = append(createdRecipes, recipe)
		createdRecipeIDs = append(createdRecipeIDs, recipe.ID)
	}

	exampleMeal := fakes.BuildFakeMeal()
	exampleMealInput := fakes.BuildFakeMealCreationRequestInputFromMeal(exampleMeal)
	exampleMealInput.Recipes = createdRecipeIDs

	createdMeal, err := client.CreateMeal(ctx, exampleMealInput)
	mustnt(err, "creating meal")

	return createdMeal
}

func main() {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()

	ctx := context.Background()

	logger := zerolog.NewZerologLogger()
	householdLeader, householdLeaderClient := getClientForUser(ctx)

	// create household members
	logger.Debug("determining household ID")
	currentStatus, statusErr := householdLeaderClient.UserStatus(ctx)
	mustnt(statusErr, "checking household leader user status")

	relevantHouseholdID := currentStatus.ActiveHousehold

	createdUsers := []*types.User{}
	createdClients := []*httpclient.Client{}

	for i := 0; i < 2; i++ {
		logger.Debug("creating user to invite")
		u, c := getClientForUser(ctx)

		logger.Debug("inviting user")
		invitation, err := householdLeaderClient.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
			FromUser:             householdLeader.ID,
			Note:                 "prober testing",
			ToEmail:              u.EmailAddress,
			DestinationHousehold: relevantHouseholdID,
		})
		mustnt(err, "")

		logger.Debug("checking for sent invitation")
		sentInvitations, err := householdLeaderClient.GetPendingHouseholdInvitationsFromUser(ctx, nil)
		mustnt(err, "")
		if len(sentInvitations.HouseholdInvitations) == 0 {
			panic("no sent invitations")
		}

		logger.Debug("checking for received invitation")
		invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
		mustnt(err, "")
		if len(invitations.HouseholdInvitations) == 0 {
			panic("no received invitations")
		}

		mustnt(c.AcceptHouseholdInvitation(ctx, relevantHouseholdID, invitation.ID, "prober testing"), "accepting household invitation")

		mustnt(c.SwitchActiveHousehold(ctx, relevantHouseholdID), "switching household")

		createdUsers = append(createdUsers, u)
		createdClients = append(createdClients, c)
	}

	// create recipes for meal plan
	createdMeals := []*types.Meal{}
	for i := 0; i < 3; i++ {
		createdMeal := createMealForTest(ctx, logger, householdLeaderClient)
		createdMeals = append(createdMeals, createdMeal)
	}

	exampleMealPlan := &types.MealPlan{
		Notes:          "prober testing",
		Status:         types.AwaitingVotesMealPlanStatus,
		StartsAt:       uint64(time.Now().Add(24 * time.Hour).Unix()),
		EndsAt:         uint64(time.Now().Add(72 * time.Hour).Unix()),
		VotingDeadline: uint64(time.Now().Add(10 * time.Minute).Unix()),
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
	exampleMealPlanInput.VotingDeadline = uint64(time.Now().Add(creationDeadline).Unix())
	createdMealPlan, err := householdLeaderClient.CreateMealPlan(ctx, exampleMealPlanInput)
	mustnt(err, "")

	createdMealPlan, err = householdLeaderClient.GetMealPlan(ctx, createdMealPlan.ID)
	mustnt(err, "")

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
		exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(vote)
		_, err = createdClients[0].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)
		mustnt(err, "")
	}

	for _, vote := range userBVotes {
		exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(vote)
		_, err = createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)
		mustnt(err, "")
	}

	createdMealPlan, err = householdLeaderClient.GetMealPlan(ctx, createdMealPlan.ID)
	mustnt(err, "")
	if types.AwaitingVotesMealPlanStatus != createdMealPlan.Status {
		panic("unexpected meal plan status")
	}

	time.Sleep(creationDeadline)

	createdMealPlan, err = householdLeaderClient.GetMealPlan(ctx, createdMealPlan.ID)
	mustnt(err, "")
	if types.FinalizedMealPlanStatus != createdMealPlan.Status {
		panic("unexpected final meal plan status")
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
					panic("selection wasn't made")
				}
			}
		}
	}

	logger.Info("done")
}
