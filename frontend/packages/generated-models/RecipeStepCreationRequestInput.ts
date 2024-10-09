// GENERATED CODE, DO NOT EDIT MANUALLY

 import { RecipeStepCompletionConditionCreationRequestInput } from './RecipeStepCompletionConditionCreationRequestInput';
 import { RecipeStepIngredientCreationRequestInput } from './RecipeStepIngredientCreationRequestInput';
 import { RecipeStepInstrumentCreationRequestInput } from './RecipeStepInstrumentCreationRequestInput';
 import { RecipeStepProductCreationRequestInput } from './RecipeStepProductCreationRequestInput';
 import { RecipeStepVesselCreationRequestInput } from './RecipeStepVesselCreationRequestInput';
 import { NumberRange } from './number_range';


export interface IRecipeStepCreationRequestInput {
   explicitInstructions: string;
 ingredients: RecipeStepIngredientCreationRequestInput;
 products: RecipeStepProductCreationRequestInput;
 vessels: RecipeStepVesselCreationRequestInput;
 completionConditions: RecipeStepCompletionConditionCreationRequestInput;
 startTimerAutomatically: boolean;
 temperatureInCelsius: NumberRange;
 optional: boolean;
 instruments: RecipeStepInstrumentCreationRequestInput;
 preparationID: string;
 estimatedTimeInSeconds: NumberRange;
 index: number;
 notes: string;
 conditionExpression: string;

}

export class RecipeStepCreationRequestInput implements IRecipeStepCreationRequestInput {
   explicitInstructions: string;
 ingredients: RecipeStepIngredientCreationRequestInput;
 products: RecipeStepProductCreationRequestInput;
 vessels: RecipeStepVesselCreationRequestInput;
 completionConditions: RecipeStepCompletionConditionCreationRequestInput;
 startTimerAutomatically: boolean;
 temperatureInCelsius: NumberRange;
 optional: boolean;
 instruments: RecipeStepInstrumentCreationRequestInput;
 preparationID: string;
 estimatedTimeInSeconds: NumberRange;
 index: number;
 notes: string;
 conditionExpression: string;
constructor(input: Partial<RecipeStepCreationRequestInput> = {}) {
	 this.explicitInstructions = input.explicitInstructions = '';
 this.ingredients = input.ingredients = new RecipeStepIngredientCreationRequestInput();
 this.products = input.products = new RecipeStepProductCreationRequestInput();
 this.vessels = input.vessels = new RecipeStepVesselCreationRequestInput();
 this.completionConditions = input.completionConditions = new RecipeStepCompletionConditionCreationRequestInput();
 this.startTimerAutomatically = input.startTimerAutomatically = false;
 this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
 this.optional = input.optional = false;
 this.instruments = input.instruments = new RecipeStepInstrumentCreationRequestInput();
 this.preparationID = input.preparationID = '';
 this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
 this.index = input.index = 0;
 this.notes = input.notes = '';
 this.conditionExpression = input.conditionExpression = '';
}
}