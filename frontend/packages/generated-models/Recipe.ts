// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeMedia } from './RecipeMedia';
import { RecipePrepTask } from './RecipePrepTask';
import { RecipeStep } from './RecipeStep';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipe {
  createdAt: string;
  eligibleForMeals: boolean;
  portionName: string;
  name: string;
  prepTasks: RecipePrepTask;
  archivedAt?: string;
  createdByUser: string;
  description: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  id: string;
  inspiredByRecipeID?: string;
  sealOfApproval: boolean;
  source: string;
  steps: RecipeStep;
  supportingRecipes: Recipe;
  lastUpdatedAt?: string;
  media: RecipeMedia;
  pluralPortionName: string;
  slug: string;
  yieldsComponentType: string;
}

export class Recipe implements IRecipe {
  createdAt: string;
  eligibleForMeals: boolean;
  portionName: string;
  name: string;
  prepTasks: RecipePrepTask;
  archivedAt?: string;
  createdByUser: string;
  description: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  id: string;
  inspiredByRecipeID?: string;
  sealOfApproval: boolean;
  source: string;
  steps: RecipeStep;
  supportingRecipes: Recipe;
  lastUpdatedAt?: string;
  media: RecipeMedia;
  pluralPortionName: string;
  slug: string;
  yieldsComponentType: string;
  constructor(input: Partial<Recipe> = {}) {
    this.createdAt = input.createdAt = '';
    this.eligibleForMeals = input.eligibleForMeals = false;
    this.portionName = input.portionName = '';
    this.name = input.name = '';
    this.prepTasks = input.prepTasks = new RecipePrepTask();
    this.archivedAt = input.archivedAt;
    this.createdByUser = input.createdByUser = '';
    this.description = input.description = '';
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.id = input.id = '';
    this.inspiredByRecipeID = input.inspiredByRecipeID;
    this.sealOfApproval = input.sealOfApproval = false;
    this.source = input.source = '';
    this.steps = input.steps = new RecipeStep();
    this.supportingRecipes = input.supportingRecipes = new Recipe();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.media = input.media = new RecipeMedia();
    this.pluralPortionName = input.pluralPortionName = '';
    this.slug = input.slug = '';
    this.yieldsComponentType = input.yieldsComponentType = '';
  }
}
