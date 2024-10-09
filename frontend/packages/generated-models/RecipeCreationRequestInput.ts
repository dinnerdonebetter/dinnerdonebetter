// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipePrepTaskWithinRecipeCreationRequestInput } from './RecipePrepTaskWithinRecipeCreationRequestInput';
 import { RecipeStepCreationRequestInput } from './RecipeStepCreationRequestInput';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IRecipeCreationRequestInput {
   alsoCreateMeal: boolean;
 estimatedPortions: NumberRangeWithOptionalMax;
 name: string;
 sealOfApproval: boolean;
 description: string;
 pluralPortionName: string;
 eligibleForMeals: boolean;
 portionName: string;
 source: string;
 steps: RecipeStepCreationRequestInput;
 inspiredByRecipeID?: string;
 prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput;
 slug: string;
 yieldsComponentType: string;

}

export class RecipeCreationRequestInput implements IRecipeCreationRequestInput {
   alsoCreateMeal: boolean;
 estimatedPortions: NumberRangeWithOptionalMax;
 name: string;
 sealOfApproval: boolean;
 description: string;
 pluralPortionName: string;
 eligibleForMeals: boolean;
 portionName: string;
 source: string;
 steps: RecipeStepCreationRequestInput;
 inspiredByRecipeID?: string;
 prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput;
 slug: string;
 yieldsComponentType: string;
constructor(input: Partial<RecipeCreationRequestInput> = {}) {
	 this.alsoCreateMeal = input.alsoCreateMeal = false;
 this.estimatedPortions = input.estimatedPortions = { min: 0 };
 this.name = input.name = '';
 this.sealOfApproval = input.sealOfApproval = false;
 this.description = input.description = '';
 this.pluralPortionName = input.pluralPortionName = '';
 this.eligibleForMeals = input.eligibleForMeals = false;
 this.portionName = input.portionName = '';
 this.source = input.source = '';
 this.steps = input.steps = new RecipeStepCreationRequestInput();
 this.inspiredByRecipeID = input.inspiredByRecipeID;
 this.prepTasks = input.prepTasks = new RecipePrepTaskWithinRecipeCreationRequestInput();
 this.slug = input.slug = '';
 this.yieldsComponentType = input.yieldsComponentType = '';
}
}