// GENERATED CODE, DO NOT EDIT MANUALLY

 import { OptionalNumberRange } from './number_range';


export interface IRecipeStepInstrumentUpdateRequestInput {
   belongsToRecipeStep: string;
 instrumentID: string;
 name: string;
 notes: string;
 optionIndex: number;
 optional: boolean;
 preferenceRank: number;
 quantity: OptionalNumberRange;
 recipeStepProductID: string;

}

export class RecipeStepInstrumentUpdateRequestInput implements IRecipeStepInstrumentUpdateRequestInput {
   belongsToRecipeStep: string;
 instrumentID: string;
 name: string;
 notes: string;
 optionIndex: number;
 optional: boolean;
 preferenceRank: number;
 quantity: OptionalNumberRange;
 recipeStepProductID: string;
constructor(input: Partial<RecipeStepInstrumentUpdateRequestInput> = {}) {
	 this.belongsToRecipeStep = input.belongsToRecipeStep || '';
 this.instrumentID = input.instrumentID || '';
 this.name = input.name || '';
 this.notes = input.notes || '';
 this.optionIndex = input.optionIndex || 0;
 this.optional = input.optional || false;
 this.preferenceRank = input.preferenceRank || 0;
 this.quantity = input.quantity || {};
 this.recipeStepProductID = input.recipeStepProductID || '';
}
}