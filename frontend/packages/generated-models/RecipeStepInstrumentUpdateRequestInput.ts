// GENERATED CODE, DO NOT EDIT MANUALLY

 import { OptionalNumberRange } from './number_range';


export interface IRecipeStepInstrumentUpdateRequestInput {
   quantity: OptionalNumberRange;
 belongsToRecipeStep?: string;
 instrumentID?: string;
 name?: string;
 optionIndex?: number;
 preferenceRank?: number;
 notes?: string;
 optional?: boolean;
 recipeStepProductID?: string;

}

export class RecipeStepInstrumentUpdateRequestInput implements IRecipeStepInstrumentUpdateRequestInput {
   quantity: OptionalNumberRange;
 belongsToRecipeStep?: string;
 instrumentID?: string;
 name?: string;
 optionIndex?: number;
 preferenceRank?: number;
 notes?: string;
 optional?: boolean;
 recipeStepProductID?: string;
constructor(input: Partial<RecipeStepInstrumentUpdateRequestInput> = {}) {
	 this.quantity = input.quantity = {};
 this.belongsToRecipeStep = input.belongsToRecipeStep;
 this.instrumentID = input.instrumentID;
 this.name = input.name;
 this.optionIndex = input.optionIndex;
 this.preferenceRank = input.preferenceRank;
 this.notes = input.notes;
 this.optional = input.optional;
 this.recipeStepProductID = input.recipeStepProductID;
}
}