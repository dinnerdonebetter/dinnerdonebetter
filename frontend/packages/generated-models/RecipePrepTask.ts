// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipePrepTaskStep } from './RecipePrepTaskStep';
 import { NumberRangeWithOptionalMax, NumberRange } from './number_range';


export interface IRecipePrepTask {
   timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
 notes: string;
 storageType: string;
 archivedAt?: string;
 belongsToRecipe: string;
 lastUpdatedAt?: string;
 name: string;
 recipeSteps: RecipePrepTaskStep;
 storageTemperatureInCelsius: NumberRange;
 description: string;
 explicitStorageInstructions: string;
 optional: boolean;
 createdAt: string;
 id: string;

}

export class RecipePrepTask implements IRecipePrepTask {
   timeBufferBeforeRecipeInSeconds: NumberRangeWithOptionalMax;
 notes: string;
 storageType: string;
 archivedAt?: string;
 belongsToRecipe: string;
 lastUpdatedAt?: string;
 name: string;
 recipeSteps: RecipePrepTaskStep;
 storageTemperatureInCelsius: NumberRange;
 description: string;
 explicitStorageInstructions: string;
 optional: boolean;
 createdAt: string;
 id: string;
constructor(input: Partial<RecipePrepTask> = {}) {
	 this.timeBufferBeforeRecipeInSeconds = input.timeBufferBeforeRecipeInSeconds = { min: 0 };
 this.notes = input.notes = '';
 this.storageType = input.storageType = '';
 this.archivedAt = input.archivedAt;
 this.belongsToRecipe = input.belongsToRecipe = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.name = input.name = '';
 this.recipeSteps = input.recipeSteps = new RecipePrepTaskStep();
 this.storageTemperatureInCelsius = input.storageTemperatureInCelsius = { min: 0, max: 0 };
 this.description = input.description = '';
 this.explicitStorageInstructions = input.explicitStorageInstructions = '';
 this.optional = input.optional = false;
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
}
}