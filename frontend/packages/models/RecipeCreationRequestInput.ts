// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskWithinRecipeCreationRequestInput } from './RecipePrepTaskWithinRecipeCreationRequestInput';
import { RecipeStepCreationRequestInput } from './RecipeStepCreationRequestInput';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeCreationRequestInput {
  inspiredByRecipeID: string;
  portionName: string;
  source: string;
  name: string;
  pluralPortionName: string;
  sealOfApproval: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  steps: RecipeStepCreationRequestInput[];
  alsoCreateMeal: boolean;
  description: string;
  eligibleForMeals: boolean;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput[];
  slug: string;
  yieldsComponentType: string;
}

export class RecipeCreationRequestInput implements IRecipeCreationRequestInput {
  inspiredByRecipeID: string;
  portionName: string;
  source: string;
  name: string;
  pluralPortionName: string;
  sealOfApproval: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  steps: RecipeStepCreationRequestInput[];
  alsoCreateMeal: boolean;
  description: string;
  eligibleForMeals: boolean;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput[];
  slug: string;
  yieldsComponentType: string;
  constructor(input: Partial<RecipeCreationRequestInput> = {}) {
    this.inspiredByRecipeID = input.inspiredByRecipeID || '';
    this.portionName = input.portionName || '';
    this.source = input.source || '';
    this.name = input.name || '';
    this.pluralPortionName = input.pluralPortionName || '';
    this.sealOfApproval = input.sealOfApproval || false;
    this.estimatedPortions = input.estimatedPortions || { min: 0 };
    this.steps = input.steps || [];
    this.alsoCreateMeal = input.alsoCreateMeal || false;
    this.description = input.description || '';
    this.eligibleForMeals = input.eligibleForMeals || false;
    this.prepTasks = input.prepTasks || [];
    this.slug = input.slug || '';
    this.yieldsComponentType = input.yieldsComponentType || '';
  }
}
