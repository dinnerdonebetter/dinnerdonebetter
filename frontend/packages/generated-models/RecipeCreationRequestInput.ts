// GENERATED CODE, DO NOT EDIT MANUALLY

import { RecipePrepTaskWithinRecipeCreationRequestInput } from './RecipePrepTaskWithinRecipeCreationRequestInput';
import { RecipeStepCreationRequestInput } from './RecipeStepCreationRequestInput';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeCreationRequestInput {
  alsoCreateMeal: boolean;
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  sealOfApproval: boolean;
  name: string;
  pluralPortionName: string;
  source: string;
  inspiredByRecipeID?: string;
  portionName: string;
  slug: string;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput;
  steps: RecipeStepCreationRequestInput;
  yieldsComponentType: string;
}

export class RecipeCreationRequestInput implements IRecipeCreationRequestInput {
  alsoCreateMeal: boolean;
  description: string;
  eligibleForMeals: boolean;
  estimatedPortions: NumberRangeWithOptionalMax;
  sealOfApproval: boolean;
  name: string;
  pluralPortionName: string;
  source: string;
  inspiredByRecipeID?: string;
  portionName: string;
  slug: string;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput;
  steps: RecipeStepCreationRequestInput;
  yieldsComponentType: string;
  constructor(input: Partial<RecipeCreationRequestInput> = {}) {
    this.alsoCreateMeal = input.alsoCreateMeal = false;
    this.description = input.description = '';
    this.eligibleForMeals = input.eligibleForMeals = false;
    this.estimatedPortions = input.estimatedPortions = { min: 0 };
    this.sealOfApproval = input.sealOfApproval = false;
    this.name = input.name = '';
    this.pluralPortionName = input.pluralPortionName = '';
    this.source = input.source = '';
    this.inspiredByRecipeID = input.inspiredByRecipeID;
    this.portionName = input.portionName = '';
    this.slug = input.slug = '';
    this.prepTasks = input.prepTasks = new RecipePrepTaskWithinRecipeCreationRequestInput();
    this.steps = input.steps = new RecipeStepCreationRequestInput();
    this.yieldsComponentType = input.yieldsComponentType = '';
  }
}
