// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskWithinRecipeCreationRequestInput } from './RecipePrepTaskWithinRecipeCreationRequestInput';
import { RecipeStepCreationRequestInput } from './RecipeStepCreationRequestInput';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeCreationRequestInput {
  sealOfApproval: boolean;
  steps: RecipeStepCreationRequestInput;
  pluralPortionName: string;
  portionName: string;
  slug: string;
  yieldsComponentType: string;
  inspiredByRecipeID?: string;
  name: string;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput;
  source: string;
  alsoCreateMeal: boolean;
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
}

export class RecipeCreationRequestInput implements IRecipeCreationRequestInput {
  sealOfApproval: boolean;
  steps: RecipeStepCreationRequestInput;
  pluralPortionName: string;
  portionName: string;
  slug: string;
  yieldsComponentType: string;
  inspiredByRecipeID?: string;
  name: string;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput;
  source: string;
  alsoCreateMeal: boolean;
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  constructor(input: Partial<RecipeCreationRequestInput> = {}) {
    this.sealOfApproval = input.sealOfApproval = false;
    this.steps = input.steps = new RecipeStepCreationRequestInput();
    this.pluralPortionName = input.pluralPortionName = '';
    this.portionName = input.portionName = '';
    this.slug = input.slug = '';
    this.yieldsComponentType = input.yieldsComponentType = '';
    this.inspiredByRecipeID = input.inspiredByRecipeID;
    this.name = input.name = '';
    this.prepTasks = input.prepTasks = new RecipePrepTaskWithinRecipeCreationRequestInput();
    this.source = input.source = '';
    this.alsoCreateMeal = input.alsoCreateMeal = false;
    this.description = input.description = '';
    this.eligibleForMeals = input.eligibleForMeals = false;
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
  }
}
