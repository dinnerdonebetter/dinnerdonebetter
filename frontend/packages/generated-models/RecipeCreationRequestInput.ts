// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskWithinRecipeCreationRequestInput } from './RecipePrepTaskWithinRecipeCreationRequestInput';
import { RecipeStepCreationRequestInput } from './RecipeStepCreationRequestInput';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeCreationRequestInput {
  eligibleForMeals: boolean;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput;
  alsoCreateMeal: boolean;
  pluralPortionName: string;
  sealOfApproval: boolean;
  slug: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  inspiredByRecipeID?: string;
  name: string;
  portionName: string;
  steps: RecipeStepCreationRequestInput;
  description: string;
  source: string;
  yieldsComponentType: string;
}

export class RecipeCreationRequestInput implements IRecipeCreationRequestInput {
  eligibleForMeals: boolean;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput;
  alsoCreateMeal: boolean;
  pluralPortionName: string;
  sealOfApproval: boolean;
  slug: string;
  estimatedPortions: NumberRangeWithOptionalMax;
  inspiredByRecipeID?: string;
  name: string;
  portionName: string;
  steps: RecipeStepCreationRequestInput;
  description: string;
  source: string;
  yieldsComponentType: string;
  constructor(input: Partial<RecipeCreationRequestInput> = {}) {
    this.eligibleForMeals = input.eligibleForMeals = false;
    this.prepTasks = input.prepTasks = new RecipePrepTaskWithinRecipeCreationRequestInput();
    this.alsoCreateMeal = input.alsoCreateMeal = false;
    this.pluralPortionName = input.pluralPortionName = '';
    this.sealOfApproval = input.sealOfApproval = false;
    this.slug = input.slug = '';
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.inspiredByRecipeID = input.inspiredByRecipeID;
    this.name = input.name = '';
    this.portionName = input.portionName = '';
    this.steps = input.steps = new RecipeStepCreationRequestInput();
    this.description = input.description = '';
    this.source = input.source = '';
    this.yieldsComponentType = input.yieldsComponentType = '';
  }
}
