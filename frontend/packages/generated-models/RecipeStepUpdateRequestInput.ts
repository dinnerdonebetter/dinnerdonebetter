// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidPreparation } from './ValidPreparation';
 import { NumberRange } from './number_range';


export interface IRecipeStepUpdateRequestInput {
   estimatedTimeInSeconds: NumberRange;
 explicitInstructions?: string;
 optional?: boolean;
 temperatureInCelsius: NumberRange;
 conditionExpression?: string;
 index?: number;
 notes?: string;
 preparation?: ValidPreparation;
 startTimerAutomatically?: boolean;
 belongsToRecipe: string;

}

export class RecipeStepUpdateRequestInput implements IRecipeStepUpdateRequestInput {
   estimatedTimeInSeconds: NumberRange;
 explicitInstructions?: string;
 optional?: boolean;
 temperatureInCelsius: NumberRange;
 conditionExpression?: string;
 index?: number;
 notes?: string;
 preparation?: ValidPreparation;
 startTimerAutomatically?: boolean;
 belongsToRecipe: string;
constructor(input: Partial<RecipeStepUpdateRequestInput> = {}) {
	 this.estimatedTimeInSeconds = input.estimatedTimeInSeconds = { min: 0, max: 0 };
 this.explicitInstructions = input.explicitInstructions;
 this.optional = input.optional;
 this.temperatureInCelsius = input.temperatureInCelsius = { min: 0, max: 0 };
 this.conditionExpression = input.conditionExpression;
 this.index = input.index;
 this.notes = input.notes;
 this.preparation = input.preparation;
 this.startTimerAutomatically = input.startTimerAutomatically;
 this.belongsToRecipe = input.belongsToRecipe = '';
}
}