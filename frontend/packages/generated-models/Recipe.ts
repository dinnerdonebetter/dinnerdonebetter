// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeMedia } from './RecipeMedia';
import { RecipePrepTask } from './RecipePrepTask';
import { RecipeStep } from './RecipeStep';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipe {
  yieldsComponentType: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  inspiredByRecipeID?: string;
  lastUpdatedAt?: string;
  name: string;
  source: string;
  createdAt: string;
  createdByUser: string;
  description: string;
  pluralPortionName: string;
  portionName: string;
  steps: RecipeStep;
  id: string;
  media: RecipeMedia;
  sealOfApproval: boolean;
  archivedAt?: string;
  prepTasks: RecipePrepTask;
  slug: string;
  supportingRecipes: Recipe;
}

export class Recipe implements IRecipe {
  yieldsComponentType: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  inspiredByRecipeID?: string;
  lastUpdatedAt?: string;
  name: string;
  source: string;
  createdAt: string;
  createdByUser: string;
  description: string;
  pluralPortionName: string;
  portionName: string;
  steps: RecipeStep;
  id: string;
  media: RecipeMedia;
  sealOfApproval: boolean;
  archivedAt?: string;
  prepTasks: RecipePrepTask;
  slug: string;
  supportingRecipes: Recipe;
  constructor(input: Partial<Recipe> = {}) {
    this.yieldsComponentType = input.yieldsComponentType = '';
    this.eligibleForMeals = input.eligibleForMeals = false;
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.inspiredByRecipeID = input.inspiredByRecipeID;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.source = input.source = '';
    this.createdAt = input.createdAt = '';
    this.createdByUser = input.createdByUser = '';
    this.description = input.description = '';
    this.pluralPortionName = input.pluralPortionName = '';
    this.portionName = input.portionName = '';
    this.steps = input.steps = new RecipeStep();
    this.id = input.id = '';
    this.media = input.media = new RecipeMedia();
    this.sealOfApproval = input.sealOfApproval = false;
    this.archivedAt = input.archivedAt;
    this.prepTasks = input.prepTasks = new RecipePrepTask();
    this.slug = input.slug = '';
    this.supportingRecipes = input.supportingRecipes = new Recipe();
  }
}
