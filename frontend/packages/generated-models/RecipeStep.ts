// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipeMedia } from './RecipeMedia';
 import { RecipeStepCompletionCondition } from './RecipeStepCompletionCondition';
 import { RecipeStepIngredient } from './RecipeStepIngredient';
 import { RecipeStepInstrument } from './RecipeStepInstrument';
 import { RecipeStepProduct } from './RecipeStepProduct';
 import { RecipeStepVessel } from './RecipeStepVessel';
 import { ValidPreparation } from './ValidPreparation';
 import { NumberRange } from './number_range';


export interface IRecipeStep {
   archivedAt?: string;
 estimatedTimeInSeconds: NumberRange;
 ingredients: RecipeStepIngredient;
 optional: boolean;
 preparation: ValidPreparation;
 conditionExpression: string;
 explicitInstructions: string;
 id: string;
 index: number;
 media: RecipeMedia;
 startTimerAutomatically: boolean;
 temperatureInCelsius: NumberRange;
 belongsToRecipe: string;
 completionConditions: RecipeStepCompletionCondition;
 instruments: RecipeStepInstrument;
 products: RecipeStepProduct;
 createdAt: string;
 lastUpdatedAt?: string;
 notes: string;
 vessels: RecipeStepVessel;

}

export class RecipeStep implements IRecipeStep {
   archivedAt?: string;
 estimatedTimeInSeconds: NumberRange;
 ingredients: RecipeStepIngredient;
 optional: boolean;
 preparation: ValidPreparation;
 conditionExpression: string;
 explicitInstructions: string;
 id: string;
 index: number;
 media: RecipeMedia;
 startTimerAutomatically: boolean;
 temperatureInCelsius: NumberRange;
 belongsToRecipe: string;
 completionConditions: RecipeStepCompletionCondition;
 instruments: RecipeStepInstrument;
 products: RecipeStepProduct;
 createdAt: string;
 lastUpdatedAt?: string;
 notes: string;
 vessels: RecipeStepVessel;
constructor(input: Partial<RecipeStep> = {}) {
	 this.archivedAt = input.archivedAt;
 this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
 this.ingredients = input.ingredients = new RecipeStepIngredient();
 this.optional = input.optional = false;
 this.preparation = input.preparation = new ValidPreparation();
 this.conditionExpression = input.conditionExpression = '';
 this.explicitInstructions = input.explicitInstructions = '';
 this.id = input.id = '';
 this.index = input.index = 0;
 this.media = input.media = new RecipeMedia();
 this.startTimerAutomatically = input.startTimerAutomatically = false;
 this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
 this.belongsToRecipe = input.belongsToRecipe = '';
 this.completionConditions = input.completionConditions = new RecipeStepCompletionCondition();
 this.instruments = input.instruments = new RecipeStepInstrument();
 this.products = input.products = new RecipeStepProduct();
 this.createdAt = input.createdAt = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.notes = input.notes = '';
 this.vessels = input.vessels = new RecipeStepVessel();
}
}