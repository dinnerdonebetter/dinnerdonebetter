// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipeMedia } from './RecipeMedia';
import { RecipePrepTask } from './RecipePrepTask';
import { RecipeStep } from './RecipeStep';
import { NumberRangeWithOptionalMax } from './number_range.gen';

export interface IRecipe {
  archivedAt: string;
  createdAt: string;
  createdByUser: string;
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  id: string;
  inspiredByRecipeID: string;
  lastUpdatedAt: string;
  media: RecipeMedia[];
  name: string;
  pluralPortionName: string;
  portionName: string;
  prepTasks: RecipePrepTask[];
  sealOfApproval: boolean;
  slug: string;
  source: string;
  steps: RecipeStep[];
  supportingRecipes: Recipe[];
  yieldsComponentType: string;
}

export class Recipe implements IRecipe {
  archivedAt: string;
  createdAt: string;
  createdByUser: string;
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  id: string;
  inspiredByRecipeID: string;
  lastUpdatedAt: string;
  media: RecipeMedia[];
  name: string;
  pluralPortionName: string;
  portionName: string;
  prepTasks: RecipePrepTask[];
  sealOfApproval: boolean;
  slug: string;
  source: string;
  steps: RecipeStep[];
  supportingRecipes: Recipe[];
  yieldsComponentType: string;
  constructor(input: Partial<Recipe> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.createdByUser = input.createdByUser || '';
    this.description = input.description || '';
    this.eligibleForMeals = input.eligibleForMeals || false;
    this.estimatedPortions = input.estimatedPortions || { min: 0 };
    this.id = input.id || '';
    this.inspiredByRecipeID = input.inspiredByRecipeID || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.media = input.media || [];
    this.name = input.name || '';
    this.pluralPortionName = input.pluralPortionName || '';
    this.portionName = input.portionName || '';
    this.prepTasks = input.prepTasks || [];
    this.sealOfApproval = input.sealOfApproval || false;
    this.slug = input.slug || '';
    this.source = input.source || '';
    this.steps = input.steps || [];
    this.supportingRecipes = input.supportingRecipes || [];
    this.yieldsComponentType = input.yieldsComponentType || '';
  }
}
