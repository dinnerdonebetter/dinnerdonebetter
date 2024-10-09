// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeMedia } from './RecipeMedia';
import { RecipePrepTask } from './RecipePrepTask';
import { RecipeStep } from './RecipeStep';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipe {
  inspiredByRecipeID?: string;
  source: string;
  yieldsComponentType: string;
  createdByUser: string;
  media: RecipeMedia;
  prepTasks: RecipePrepTask;
  slug: string;
  steps: RecipeStep;
  eligibleForMeals: boolean;
  id: string;
  pluralPortionName: string;
  supportingRecipes: Recipe;
  portionName: string;
  sealOfApproval: boolean;
  archivedAt?: string;
  createdAt: string;
  description: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  lastUpdatedAt?: string;
  name: string;
}

export class Recipe implements IRecipe {
  inspiredByRecipeID?: string;
  source: string;
  yieldsComponentType: string;
  createdByUser: string;
  media: RecipeMedia;
  prepTasks: RecipePrepTask;
  slug: string;
  steps: RecipeStep;
  eligibleForMeals: boolean;
  id: string;
  pluralPortionName: string;
  supportingRecipes: Recipe;
  portionName: string;
  sealOfApproval: boolean;
  archivedAt?: string;
  createdAt: string;
  description: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  lastUpdatedAt?: string;
  name: string;
  constructor(input: Partial<Recipe> = {}) {
    this.inspiredByRecipeID = input.inspiredByRecipeID;
    this.source = input.source = '';
    this.yieldsComponentType = input.yieldsComponentType = '';
    this.createdByUser = input.createdByUser = '';
    this.media = input.media = new RecipeMedia();
    this.prepTasks = input.prepTasks = new RecipePrepTask();
    this.slug = input.slug = '';
    this.steps = input.steps = new RecipeStep();
    this.eligibleForMeals = input.eligibleForMeals = false;
    this.id = input.id = '';
    this.pluralPortionName = input.pluralPortionName = '';
    this.supportingRecipes = input.supportingRecipes = new Recipe();
    this.portionName = input.portionName = '';
    this.sealOfApproval = input.sealOfApproval = false;
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
  }
}
