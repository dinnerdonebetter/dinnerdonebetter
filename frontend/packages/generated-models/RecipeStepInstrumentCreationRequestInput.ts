// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRangeWithOptionalMax } from './number_range';


export interface IRecipeStepInstrumentCreationRequestInput {
   optional: boolean;
 productOfRecipeStepProductIndex?: number;
 instrumentID?: string;
 name: string;
 notes: string;
 quantity: NumberRangeWithOptionalMax;
 recipeStepProductID?: string;
 optionIndex: number;
 preferenceRank: number;
 productOfRecipeStepIndex?: number;

}

export class RecipeStepInstrumentCreationRequestInput implements IRecipeStepInstrumentCreationRequestInput {
   optional: boolean;
 productOfRecipeStepProductIndex?: number;
 instrumentID?: string;
 name: string;
 notes: string;
 quantity: NumberRangeWithOptionalMax;
 recipeStepProductID?: string;
 optionIndex: number;
 preferenceRank: number;
 productOfRecipeStepIndex?: number;
constructor(input: Partial<RecipeStepInstrumentCreationRequestInput> = {}) {
	 this.optional = input.optional = false;
 this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
 this.instrumentID = input.instrumentID;
 this.name = input.name = '';
 this.notes = input.notes = '';
 this.quantity = input.quantity = { min: 0 };
 this.recipeStepProductID = input.recipeStepProductID;
 this.optionIndex = input.optionIndex = 0;
 this.preferenceRank = input.preferenceRank = 0;
 this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
}
}