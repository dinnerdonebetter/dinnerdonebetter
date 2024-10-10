// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskWithinRecipeCreationRequestInput } from './RecipePrepTaskWithinRecipeCreationRequestInput';
import { RecipeStepCreationRequestInput } from './RecipeStepCreationRequestInput';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeCreationRequestInput {
  alsoCreateMeal: boolean;
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  inspiredByRecipeID: string;
  name: string;
  pluralPortionName: string;
  portionName: string;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput[];
  sealOfApproval: boolean;
  slug: string;
  source: string;
  steps: RecipeStepCreationRequestInput[];
  yieldsComponentType: string;
}

export class RecipeCreationRequestInput implements IRecipeCreationRequestInput {
  alsoCreateMeal: boolean;
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  inspiredByRecipeID: string;
  name: string;
  pluralPortionName: string;
  portionName: string;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput[];
  sealOfApproval: boolean;
  slug: string;
  source: string;
  steps: RecipeStepCreationRequestInput[];
  yieldsComponentType: string;
  constructor(input: Partial<RecipeCreationRequestInput> = {}) {
    this.alsoCreateMeal = input.alsoCreateMeal || false;
    this.description = input.description || '';
    this.eligibleForMeals = input.eligibleForMeals || false;
    this.estimatedPortions = input.estimatedPortions || { min: 0 };
    this.inspiredByRecipeID = input.inspiredByRecipeID || '';
    this.name = input.name || '';
    this.pluralPortionName = input.pluralPortionName || '';
    this.portionName = input.portionName || '';
    this.prepTasks = input.prepTasks || [];
    this.sealOfApproval = input.sealOfApproval || false;
    this.slug = input.slug || '';
    this.source = input.source || '';
    this.steps = input.steps || [];
    this.yieldsComponentType = input.yieldsComponentType || '';
  }
}
