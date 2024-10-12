// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipePrepTaskStepWithinRecipeCreationRequestInput } from './RecipePrepTaskStepWithinRecipeCreationRequestInput';
 import { NumberRange, NumberRangeWithOptionalMax } from './number_range';


export interface IRecipePrepTaskWithinRecipeCreationRequestInput {
   belongsToRecipe: string;
 description: string;
 explicitStorageInstructions: string;
 name: string;
 notes: string;
 optional: boolean;
 recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput[];
 storageTemperatureInCelsius: NumberRange;
 storageType: string;
 timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;

}

export class RecipePrepTaskWithinRecipeCreationRequestInput implements IRecipePrepTaskWithinRecipeCreationRequestInput {
   belongsToRecipe: string;
 description: string;
 explicitStorageInstructions: string;
 name: string;
 notes: string;
 optional: boolean;
 recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput[];
 storageTemperatureInCelsius: NumberRange;
 storageType: string;
 timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
constructor(input: Partial<RecipePrepTaskWithinRecipeCreationRequestInput> = {}) {
	 this.belongsToRecipe = input.belongsToRecipe || '';
 this.description = input.description || '';
 this.explicitStorageInstructions = input.explicitStorageInstructions || '';
 this.name = input.name || '';
 this.notes = input.notes || '';
 this.optional = input.optional || false;
 this.recipeSteps = input.recipeSteps || [];
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
 this.storageType = input.storageType || '';
 this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds || { min: 0 };
}
}