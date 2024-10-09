// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipePrepTaskStepCreationRequestInput } from './RecipePrepTaskStepCreationRequestInput';
 import { NumberRange, NumberRangeWithOptionalMax } from './number_range';


export interface IRecipePrepTaskCreationRequestInput {
   notes: string;
 optional: boolean;
 storageTemperatureInCelsius: NumberRange;
 timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
 belongsToRecipe: string;
 description: string;
 recipeSteps: RecipePrepTaskStepCreationRequestInput;
 storageType: string;
 explicitStorageInstructions: string;
 name: string;

}

export class RecipePrepTaskCreationRequestInput implements IRecipePrepTaskCreationRequestInput {
   notes: string;
 optional: boolean;
 storageTemperatureInCelsius: NumberRange;
 timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
 belongsToRecipe: string;
 description: string;
 recipeSteps: RecipePrepTaskStepCreationRequestInput;
 storageType: string;
 explicitStorageInstructions: string;
 name: string;
constructor(input: Partial<RecipePrepTaskCreationRequestInput> = {}) {
	 this.notes = input.notes = '';
 this.optional = input.optional = false;
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
 this.belongsToRecipe = input.belongsToRecipe = '';
 this.description = input.description = '';
 this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepCreationRequestInput();
 this.storageType = input.storageType = '';
 this.explicitStorageInstructions = input.explicitStorageInstructions = '';
 this.name = input.name = '';
}
}