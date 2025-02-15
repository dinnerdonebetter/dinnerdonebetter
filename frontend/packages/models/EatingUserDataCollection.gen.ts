// GENERATED CODE, DO NOT EDIT MANUALLY

import { Meal } from './Meal.gen';
import { Recipe } from './Recipe.gen';
import { RecipeRating } from './RecipeRating.gen';
import { UserIngredientPreference } from './UserIngredientPreference.gen';

export interface IEatingUserDataCollection {
  householdInstrumentOwnerships: object;
  mealPlans: object;
  meals: Meal[];
  recipeRatings: RecipeRating[];
  recipes: Recipe[];
  reportID: string;
  userIngredientPreferences: UserIngredientPreference[];
}

export class EatingUserDataCollection implements IEatingUserDataCollection {
  householdInstrumentOwnerships: object;
  mealPlans: object;
  meals: Meal[];
  recipeRatings: RecipeRating[];
  recipes: Recipe[];
  reportID: string;
  userIngredientPreferences: UserIngredientPreference[];
  constructor(input: Partial<EatingUserDataCollection> = {}) {
    this.householdInstrumentOwnerships = input.householdInstrumentOwnerships || {};
    this.mealPlans = input.mealPlans || {};
    this.meals = input.meals || [];
    this.recipeRatings = input.recipeRatings || [];
    this.recipes = input.recipes || [];
    this.reportID = input.reportID || '';
    this.userIngredientPreferences = input.userIngredientPreferences || [];
  }
}
