package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// AllMeals returns all bootstrap meal creation inputs.
// Each meal pairs main dishes with appropriate side dishes.
// Each meal is created with the provided userID as the creator.
func AllMeals(userID string, recipes []*mealplanning.RecipeDatabaseCreationInput) []*mealplanning.MealDatabaseCreationInput {
	// Build maps of recipes by component type and by name for easy lookup
	mainRecipes := make(map[string]*mealplanning.RecipeDatabaseCreationInput)
	sideRecipes := make(map[string]*mealplanning.RecipeDatabaseCreationInput)
	saladRecipes := make(map[string]*mealplanning.RecipeDatabaseCreationInput)

	for _, recipe := range recipes {
		switch recipe.YieldsComponentType {
		case mealplanning.MealComponentTypesMain:
			mainRecipes[recipe.Name] = recipe
		case mealplanning.MealComponentTypesSide:
			sideRecipes[recipe.Name] = recipe
		case mealplanning.MealComponentTypesSalad:
			saladRecipes[recipe.Name] = recipe
		}
	}

	var meals []*mealplanning.MealDatabaseCreationInput

	// Helper function to calculate recipe scale
	// Scale side dish to match main dish portions
	calculateScale := func(mainPortions, sidePortions float32) float32 {
		if sidePortions == 0 {
			return 1.0
		}
		return mainPortions / sidePortions
	}

	// Helper function to create a meal component
	createComponent := func(recipe *mealplanning.RecipeDatabaseCreationInput, componentType string, scale float32) *mealplanning.MealComponentDatabaseCreationInput {
		return &mealplanning.MealComponentDatabaseCreationInput{
			RecipeID:      recipe.ID,
			ComponentType: componentType,
			RecipeScale:   scale,
		}
	}

	// Helper function to get main portion range
	getMainPortions := func(recipe *mealplanning.RecipeDatabaseCreationInput) (min float32, max *float32) {
		min = recipe.EstimatedPortions.Min
		if recipe.EstimatedPortions.Max != nil {
			max = pointer.To(*recipe.EstimatedPortions.Max)
		}
		return
	}

	// Pair 1: Pan-Seared Butter-Basted Steak with Ultra-Fluffy Mashed Potatoes
	if steak, ok := mainRecipes["Pan-Seared Butter-Basted Steak"]; ok {
		if mashedPotatoes, ok := sideRecipes["Ultra-Fluffy Mashed Potatoes"]; ok {
			mainMin, mainMax := getMainPortions(steak)
			sideMin := mashedPotatoes.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(steak, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(mashedPotatoes, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Pan-Seared Steak with Mashed Potatoes",
				Description:          "A classic steak dinner with creamy mashed potatoes.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 2: Sous Vide Chicken Breast with Simple White Rice
	if chickenBreast, ok := mainRecipes["Sous Vide Chicken Breast"]; ok {
		if rice, ok := sideRecipes["Simple White Rice"]; ok {
			mainMin, mainMax := getMainPortions(chickenBreast)
			sideMin := rice.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(chickenBreast, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(rice, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Sous Vide Chicken Breast with Rice",
				Description:          "Tender chicken breast served with fluffy white rice.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 3: Perfect Roast Chicken with Caesar Roasted Broccoli
	if roastChicken, ok := mainRecipes["Perfect Roast Chicken"]; ok {
		if broccoli, ok := sideRecipes["Caesar Roasted Broccoli"]; ok {
			mainMin, mainMax := getMainPortions(roastChicken)
			sideMin := broccoli.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(roastChicken, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(broccoli, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Roast Chicken with Caesar Broccoli",
				Description:          "Classic roast chicken with flavorful Caesar-roasted broccoli.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 4: Sous Vide Pork Chops with Haricots Verts Amandine
	if porkChops, ok := mainRecipes["Sous Vide Pork Chops"]; ok {
		if greenBeans, ok := sideRecipes["Haricots Verts Amandine"]; ok {
			mainMin, mainMax := getMainPortions(porkChops)
			sideMin := greenBeans.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(porkChops, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(greenBeans, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Sous Vide Pork Chops with Green Beans",
				Description:          "Tender pork chops with French-style green beans and almonds.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 5: Classic Smashed Burgers with Mixed Green Salad
	if burgers, ok := mainRecipes["Classic Smashed Burgers"]; ok {
		if salad, ok := saladRecipes["Mixed Green Salad"]; ok {
			mainMin, mainMax := getMainPortions(burgers)
			saladMin := salad.EstimatedPortions.Min
			scale := calculateScale(mainMin, saladMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(burgers, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(salad, mealplanning.MealComponentTypesSalad, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Classic Burgers with Mixed Green Salad",
				Description:          "Juicy smash burgers served with a fresh mixed green salad.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 6: Soy Sauce–Braised Chicken Thighs with Simple White Rice
	if chickenThighs, ok := mainRecipes["Soy Sauce–Braised Chicken Thighs"]; ok {
		if rice, ok := sideRecipes["Simple White Rice"]; ok {
			mainMin, mainMax := getMainPortions(chickenThighs)
			sideMin := rice.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(chickenThighs, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(rice, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Soy Sauce Braised Chicken Thighs with Rice",
				Description:          "Flavorful braised chicken thighs served with white rice.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 7: Grilled Pork Tenderloin with Roasted Brussels Sprouts
	if porkTenderloin, ok := mainRecipes["Grilled Pork Tenderloin"]; ok {
		if brusselsSprouts, ok := sideRecipes["Roasted Brussels Sprouts"]; ok {
			mainMin, mainMax := getMainPortions(porkTenderloin)
			sideMin := brusselsSprouts.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(porkTenderloin, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(brusselsSprouts, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Grilled Pork Tenderloin with Brussels Sprouts",
				Description:          "Tender grilled pork with roasted Brussels sprouts.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 8: Crispy Pan-Seared Salmon Fillets with Caesar Roasted Broccoli
	if salmon, ok := mainRecipes["Crispy Pan-Seared Salmon Fillets"]; ok {
		if broccoli, ok := sideRecipes["Caesar Roasted Broccoli"]; ok {
			mainMin, mainMax := getMainPortions(salmon)
			sideMin := broccoli.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(salmon, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(broccoli, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Pan-Seared Salmon with Caesar Broccoli",
				Description:          "Crispy-skinned salmon with Caesar-roasted broccoli.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 9: The Best Carne Asada with Refried Beans and Tortillas
	if carneAsada, ok := mainRecipes["The Best Carne Asada"]; ok {
		components := []*mealplanning.MealComponentDatabaseCreationInput{
			createComponent(carneAsada, mealplanning.MealComponentTypesMain, 1.0),
		}

		mainMin, mainMax := getMainPortions(carneAsada)

		if refriedBeans, ok := sideRecipes["Perfect Frijoles Refritos (Mexican Refried Beans)"]; ok {
			sideMin := refriedBeans.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)
			components = append(components, createComponent(refriedBeans, mealplanning.MealComponentTypesSide, scale))
		}

		if tortillas, ok := sideRecipes["Tortillas"]; ok {
			sideMin := tortillas.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)
			components = append(components, createComponent(tortillas, mealplanning.MealComponentTypesSide, scale))
		}

		meals = append(meals, &mealplanning.MealDatabaseCreationInput{
			ID:                   identifiers.New(),
			CreatedByUser:        userID,
			Name:                 "Carne Asada with Refried Beans and Tortillas",
			Description:          "Traditional Mexican grilled beef with refried beans and fresh tortillas.",
			EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
			Components:           components,
			EligibleForMealPlans: true,
		})
	}

	// Pair 10: Butter Chicken with Simple White Rice
	if butterChicken, ok := mainRecipes["Butter Chicken"]; ok {
		if rice, ok := sideRecipes["Simple White Rice"]; ok {
			mainMin, mainMax := getMainPortions(butterChicken)
			sideMin := rice.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(butterChicken, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(rice, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Butter Chicken with Rice",
				Description:          "Creamy Indian butter chicken served with white rice.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 11: Stovetop Mac and Cheese with Caesar Salad
	if macAndCheese, ok := mainRecipes["Stovetop Mac and Cheese"]; ok {
		if caesarSalad, ok := saladRecipes["Caesar Salad"]; ok {
			mainMin, mainMax := getMainPortions(macAndCheese)
			saladMin := caesarSalad.EstimatedPortions.Min
			scale := calculateScale(mainMin, saladMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(macAndCheese, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(caesarSalad, mealplanning.MealComponentTypesSalad, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Mac and Cheese with Caesar Salad",
				Description:          "Creamy stovetop mac and cheese with classic Caesar salad.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 12: Grilled Whole Cauliflower with Teriyaki Sauce with Glazed Carrots
	if cauliflower, ok := mainRecipes["Grilled Whole Cauliflower with Teriyaki Sauce"]; ok {
		if carrots, ok := sideRecipes["Glazed Carrots with Brown Butter and Sage"]; ok {
			mainMin, mainMax := getMainPortions(cauliflower)
			sideMin := carrots.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(cauliflower, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(carrots, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Grilled Whole Cauliflower with Glazed Carrots",
				Description:          "Smoky grilled cauliflower with sweet glazed carrots.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 13: Crispy Pan-Seared Salmon with Stir-Fried Green Beans
	if salmon, ok := mainRecipes["Crispy Pan-Seared Salmon Fillets"]; ok {
		if greenBeans, ok := sideRecipes["Stir-Fried Green Beans"]; ok {
			mainMin, mainMax := getMainPortions(salmon)
			sideMin := greenBeans.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(salmon, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(greenBeans, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Pan-Seared Salmon with Stir-Fried Green Beans",
				Description:          "Crispy salmon with wok-tossed green beans.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	// Pair 14: Sous Vide Pork Chops with Cornbread
	if porkChops, ok := mainRecipes["Sous Vide Pork Chops"]; ok {
		if cornbread, ok := sideRecipes["Cornbread"]; ok {
			mainMin, mainMax := getMainPortions(porkChops)
			sideMin := cornbread.EstimatedPortions.Min
			scale := calculateScale(mainMin, sideMin)

			components := []*mealplanning.MealComponentDatabaseCreationInput{
				createComponent(porkChops, mealplanning.MealComponentTypesMain, 1.0),
				createComponent(cornbread, mealplanning.MealComponentTypesSide, scale),
			}

			meals = append(meals, &mealplanning.MealDatabaseCreationInput{
				ID:                   identifiers.New(),
				CreatedByUser:        userID,
				Name:                 "Sous Vide Pork Chops with Cornbread",
				Description:          "Tender pork chops with sweet cornbread.",
				EstimatedPortions:    types.Float32RangeWithOptionalMax{Min: mainMin, Max: mainMax},
				Components:           components,
				EligibleForMealPlans: true,
			})
		}
	}

	return meals
}
