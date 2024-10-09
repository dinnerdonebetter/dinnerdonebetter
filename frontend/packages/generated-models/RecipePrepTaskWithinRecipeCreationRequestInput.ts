// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipePrepTaskStepWithinRecipeCreationRequestInput } from './RecipePrepTaskStepWithinRecipeCreationRequestInput';
 import { NumberRangeWithOptionalMax, NumberRange } from './number_range';


export interface IRecipePrepTaskWithinRecipeCreationRequestInput {
   belongsToRecipe: string;
 explicitStorageInstructions: string;
 name: string;
 optional: boolean;
 timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
 description: string;
 notes: string;
 recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput;
 storageTemperatureInCelsius: NumberRange;
 storageType: string;

}

export class RecipePrepTaskWithinRecipeCreationRequestInput implements IRecipePrepTaskWithinRecipeCreationRequestInput {
   belongsToRecipe: string;
 explicitStorageInstructions: string;
 name: string;
 optional: boolean;
 timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
 description: string;
 notes: string;
 recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput;
 storageTemperatureInCelsius: NumberRange;
 storageType: string;
constructor(input: Partial<RecipePrepTaskWithinRecipeCreationRequestInput> = {}) {
	 this.belongsToRecipe = input.belongsToRecipe = '';
 this.explicitStorageInstructions = input.explicitStorageInstructions = '';
 this.name = input.name = '';
 this.optional = input.optional = false;
 this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
 this.description = input.description = '';
 this.notes = input.notes = '';
 this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepWithinRecipeCreationRequestInput();
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.storageType = input.storageType = '';
}
}