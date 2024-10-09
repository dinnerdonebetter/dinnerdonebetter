// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipePrepTaskStepUpdateRequestInput } from './RecipePrepTaskStepUpdateRequestInput';
 import { NumberRange, OptionalNumberRange } from './number_range';


export interface IRecipePrepTaskUpdateRequestInput {
   belongsToRecipe?: string;
 optional?: boolean;
 recipeSteps: RecipePrepTaskStepUpdateRequestInput;
 storageTemperatureInCelsius: NumberRange;
 timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
 description?: string;
 explicitStorageInstructions?: string;
 name?: string;
 notes?: string;
 storageType?: string;

}

export class RecipePrepTaskUpdateRequestInput implements IRecipePrepTaskUpdateRequestInput {
   belongsToRecipe?: string;
 optional?: boolean;
 recipeSteps: RecipePrepTaskStepUpdateRequestInput;
 storageTemperatureInCelsius: NumberRange;
 timeBufferBeforeRecipeInSeconds: OptionalNumberRange;
 description?: string;
 explicitStorageInstructions?: string;
 name?: string;
 notes?: string;
 storageType?: string;
constructor(input: Partial<RecipePrepTaskUpdateRequestInput> = {}) {
	 this.belongsToRecipe = input.belongsToRecipe;
 this.optional = input.optional;
 this.recipeSteps = input.recipeSteps = new RecipePrepTaskStepUpdateRequestInput();
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = {};
 this.description = input.description;
 this.explicitStorageInstructions = input.explicitStorageInstructions;
 this.name = input.name;
 this.notes = input.notes;
 this.storageType = input.storageType;
}
}