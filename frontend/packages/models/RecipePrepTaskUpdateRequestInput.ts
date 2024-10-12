// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipePrepTaskStepUpdateRequestInput } from './RecipePrepTaskStepUpdateRequestInput';
 import { NumberRange, OptionalNumberRange } from './number_range';


export interface IRecipePrepTaskUpdateRequestInput {
   belongsToRecipe: string;
 description: string;
 explicitStorageInstructions: string;
 name: string;
 notes: string;
 optional: boolean;
 recipeSteps: RecipePrepTaskStepUpdateRequestInput[];
 storageTemperatureInCelsius: NumberRange;
 storageType: string;
 timeBufferBeforeRecipeInSeconds: OptionalNumberRange;

}

export class RecipePrepTaskUpdateRequestInput implements IRecipePrepTaskUpdateRequestInput {
   belongsToRecipe: string;
 description: string;
 explicitStorageInstructions: string;
 name: string;
 notes: string;
 optional: boolean;
 recipeSteps: RecipePrepTaskStepUpdateRequestInput[];
 storageTemperatureInCelsius: NumberRange;
 storageType: string;
 timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
constructor(input: Partial<RecipePrepTaskUpdateRequestInput> = {}) {
	 this.belongsToRecipe = input.belongsToRecipe || '';
 this.description = input.description || '';
 this.explicitStorageInstructions = input.explicitStorageInstructions || '';
 this.name = input.name || '';
 this.notes = input.notes || '';
 this.optional = input.optional || false;
 this.recipeSteps = input.recipeSteps || [];
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius || { min: 0, max: 0 };
 this.storageType = input.storageType || '';
 this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds || {};
}
}