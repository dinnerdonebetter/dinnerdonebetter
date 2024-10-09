// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeMedia } from './RecipeMedia';
import { RecipePrepTask } from './RecipePrepTask';
import { RecipeStep } from './RecipeStep';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipe {
  archivedAt: string;
  createdAt: string;
  createdByUser: string;
  slug: string;
  source: string;
  supportingRecipes: Recipe[];
  description: string;
  eligibleForMeals: boolean;
  id: string;
  inspiredByRecipeID: string;
  media: RecipeMedia[];
  name: string;
  sealOfApproval: boolean;
  yieldsComponentType: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  lastUpdatedAt: string;
  pluralPortionName: string;
  portionName: string;
  prepTasks: RecipePrepTask[];
  steps: RecipeStep[];
}

export class Recipe implements IRecipe {
  archivedAt: string;
  createdAt: string;
  createdByUser: string;
  slug: string;
  source: string;
  supportingRecipes: Recipe[];
  description: string;
  eligibleForMeals: boolean;
  id: string;
  inspiredByRecipeID: string;
  media: RecipeMedia[];
  name: string;
  sealOfApproval: boolean;
  yieldsComponentType: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  lastUpdatedAt: string;
  pluralPortionName: string;
  portionName: string;
  prepTasks: RecipePrepTask[];
  steps: RecipeStep[];
  constructor(input: Partial<Recipe> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.createdByUser = input.createdByUser || '';
    this.slug = input.slug || '';
    this.source = input.source || '';
    this.supportingRecipes = input.supportingRecipes || [];
    this.description = input.description || '';
    this.eligibleForMeals = input.eligibleForMeals || false;
    this.id = input.id || '';
    this.inspiredByRecipeID = input.inspiredByRecipeID || '';
    this.media = input.media || [];
    this.name = input.name || '';
    this.sealOfApproval = input.sealOfApproval || false;
    this.yieldsComponentType = input.yieldsComponentType || '';
    this.estimatedPortions = input.estimatedPortions || { min: 0 };
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.pluralPortionName = input.pluralPortionName || '';
    this.portionName = input.portionName || '';
    this.prepTasks = input.prepTasks || [];
    this.steps = input.steps || [];
  }
}
