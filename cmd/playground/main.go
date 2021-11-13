package main

import (
	"time"

	"resenje.org/schulze"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func finalizeMealPlan(option *types.MealPlan) error {
	candidateMap := map[string]struct{}{}

	votesByUser := map[string]schulze.Ballot{}
	for _, o := range option.Options {
		for _, v := range o.Votes {
			if votesByUser[v.ByUser] == nil {
				votesByUser[v.ByUser] = schulze.Ballot{}
			}

			if !v.Abstain {
				votesByUser[v.ByUser][v.BelongsToMealPlanOption] = int(v.Rank)
			}

			candidateMap[v.BelongsToMealPlanOption] = struct{}{}
		}
	}

	candidates := []string{}
	for c := range candidateMap {
		candidates = append(candidates, c)
	}

	e := schulze.NewVoting(candidates...)
	for _, vote := range votesByUser {
		if err := e.Vote(vote); err != nil {
			return err
		}
	}

	winners, tie := e.Compute()
	if tie {
		println(winners)
	}

	return nil
}

func main() {
	const (
		optionA = "eggs benedict"
		optionB = "scrambled eggs"
		optionC = "buttered toast"
	)

	var (
		userID1 = fakes.BuildFakeID()
		userID2 = fakes.BuildFakeID()
		userID3 = fakes.BuildFakeID()
		userID4 = fakes.BuildFakeID()
	)

	mealPlan := &types.MealPlan{
		Options: []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    0,
						ByUser:                  userID1,
					},
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    0,
						ByUser:                  userID2,
					},
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    1,
						ByUser:                  userID3,
					},
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    2,
						ByUser:                  userID4,
					},
				},
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    0,
						ByUser:                  userID3,
					},
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    1,
						ByUser:                  userID2,
					},
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    1,
						ByUser:                  userID4,
					},
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    2,
						ByUser:                  userID1,
					},
				},
			},
			{
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionC,
						Rank:                    0,
						ByUser:                  userID4,
					},

					{
						BelongsToMealPlanOption: optionC,
						Rank:                    1,
						ByUser:                  userID1,
					},
					{
						BelongsToMealPlanOption: optionC,
						Rank:                    2,
						ByUser:                  userID2,
					},
					{
						BelongsToMealPlanOption: optionC,
						Rank:                    2,
						ByUser:                  userID3,
					},
				},
			},
		},
	}

	if err := finalizeMealPlan(mealPlan); err != nil {
		panic(err)
	}
}

/*
func buildDinnerOptions(mealPlanID string) map[string]map[string][]types.MealPlanOption {
	m := map[string]map[string][]types.MealPlanOption{
		"monday": {
			string(types.BreakfastMealName): []types.MealPlanOption{
				{
					ID:       "eggs benedict",
					Day:      1,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "scrambled eggs",
					Day:      1,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "buttered toast",
					Day:      1,
					MealName: types.BreakfastMealName,
				},
			},
			string(types.LunchMealName): []types.MealPlanOption{
				{
					ID:       "croque monsieur",
					Day:      1,
					MealName: types.LunchMealName,
				},
				{
					ID:       "baloney sandwich",
					Day:      1,
					MealName: types.LunchMealName,
				},
				{
					ID:       "mac & cheese",
					Day:      1,
					MealName: types.LunchMealName,
				},
			},
			string(types.DinnerMealName): []types.MealPlanOption{
				{
					ID:       "chicken soup",
					Day:      1,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "pizza",
					Day:      1,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "chicken & waffles",
					Day:      1,
					MealName: types.DinnerMealName,
				},
			},
		},
		"tuesday": {
			string(types.BreakfastMealName): []types.MealPlanOption{
				{
					ID:       "oatmeal",
					Day:      2,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "yogurt and berries",
					Day:      2,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "avocado toast",
					Day:      2,
					MealName: types.BreakfastMealName,
				},
			},
			string(types.LunchMealName): []types.MealPlanOption{
				{
					ID:       "ramen",
					Day:      2,
					MealName: types.LunchMealName,
				},
				{
					ID:       "hot dog",
					Day:      2,
					MealName: types.LunchMealName,
				},
				{
					ID:       "chicken tikka masala",
					Day:      2,
					MealName: types.LunchMealName,
				},
			},
			string(types.DinnerMealName): []types.MealPlanOption{
				{
					ID:       "tacos",
					Day:      2,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "fisherman's pie",
					Day:      2,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "pork roast",
					Day:      2,
					MealName: types.DinnerMealName,
				},
			},
		},
		"wednesday": {
			string(types.BreakfastMealName): []types.MealPlanOption{
				{
					ID:       "biscuits and gravy",
					Day:      3,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "cereal",
					Day:      3,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "coffee",
					Day:      3,
					MealName: types.BreakfastMealName,
				},
			},
			string(types.LunchMealName): []types.MealPlanOption{
				{
					ID:       "tuna salad",
					Day:      3,
					MealName: types.LunchMealName,
				},
				{
					ID:       "chicken sandwich",
					Day:      3,
					MealName: types.LunchMealName,
				},
				{
					ID:       "fried oysters",
					Day:      3,
					MealName: types.LunchMealName,
				},
			},
			string(types.DinnerMealName): []types.MealPlanOption{
				{
					ID:       "eggplant parmesan",
					Day:      3,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "thai curry",
					Day:      3,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "fried rice",
					Day:      3,
					MealName: types.DinnerMealName,
				},
			},
		},
		"thursday": {
			string(types.BreakfastMealName): []types.MealPlanOption{
				{
					ID:       "rice",
					Day:      4,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "poached eggs",
					Day:      4,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "biscuits and fruit",
					Day:      4,
					MealName: types.BreakfastMealName,
				},
			},
			string(types.LunchMealName): []types.MealPlanOption{
				{
					ID:       "spanish tortilla",
					Day:      4,
					MealName: types.LunchMealName,
				},
				{
					ID:       "potato and ground beef skillet",
					Day:      4,
					MealName: types.LunchMealName,
				},
				{
					ID:       "beef fajita",
					Day:      4,
					MealName: types.LunchMealName,
				},
			},
			string(types.DinnerMealName): []types.MealPlanOption{
				{
					ID:       "tlacoyo",
					Day:      4,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "pozole",
					Day:      4,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "lo mein",
					Day:      4,
					MealName: types.DinnerMealName,
				},
			},
		},
		"friday": {
			string(types.BreakfastMealName): []types.MealPlanOption{
				{
					ID:       "steak and eggs",
					Day:      5,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "grits",
					Day:      5,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "fried egg",
					Day:      5,
					MealName: types.BreakfastMealName,
				},
			},
			string(types.LunchMealName): []types.MealPlanOption{
				{
					ID:       "white bean soup",
					Day:      5,
					MealName: types.LunchMealName,
				},
				{
					ID:       "chicken enchiladas",
					Day:      5,
					MealName: types.LunchMealName,
				},
				{
					ID:       "potato salad",
					Day:      5,
					MealName: types.LunchMealName,
				},
			},
			string(types.DinnerMealName): []types.MealPlanOption{
				{
					ID:       "burrito",
					Day:      5,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "hamburger casserole",
					Day:      5,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "chicken tenders",
					Day:      5,
					MealName: types.DinnerMealName,
				},
			},
		},
		"saturday": {
			string(types.BreakfastMealName): []types.MealPlanOption{
				{
					ID:       "groats",
					Day:      6,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "overnight oats",
					Day:      6,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "berries",
					Day:      6,
					MealName: types.BreakfastMealName,
				},
			},
			string(types.LunchMealName): []types.MealPlanOption{
				{
					ID:       "salt baked fish",
					Day:      6,
					MealName: types.LunchMealName,
				},
				{
					ID:       "garlic noodles",
					Day:      6,
					MealName: types.LunchMealName,
				},
				{
					ID:       "patty melt",
					Day:      6,
					MealName: types.LunchMealName,
				},
			},
			string(types.DinnerMealName): []types.MealPlanOption{
				{
					ID:       "broccoli and cheese soup",
					Day:      6,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "rice bowl",
					Day:      6,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "veggie kebab",
					Day:      6,
					MealName: types.DinnerMealName,
				},
			},
		},
		"sunday": {
			string(types.BreakfastMealName): []types.MealPlanOption{
				{
					ID:       "soft scramble",
					Day:      0,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "scotch egg",
					Day:      0,
					MealName: types.BreakfastMealName,
				},
				{
					ID:       "fruit loops",
					Day:      0,
					MealName: types.BreakfastMealName,
				},
			},
			string(types.LunchMealName): []types.MealPlanOption{
				{
					ID:       "sub sandwich",
					Day:      0,
					MealName: types.LunchMealName,
				},
				{
					ID:       "tuna casserole",
					Day:      0,
					MealName: types.LunchMealName,
				},
				{
					ID:       "pasta salad",
					Day:      0,
					MealName: types.LunchMealName,
				},
			},
			string(types.DinnerMealName): []types.MealPlanOption{
				{
					ID:       "pasta",
					Day:      0,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "croque madame",
					Day:      0,
					MealName: types.DinnerMealName,
				},
				{
					ID:       "lasagna",
					Day:      0,
					MealName: types.DinnerMealName,
				},
			},
		},
	}
	return m
}
*/
