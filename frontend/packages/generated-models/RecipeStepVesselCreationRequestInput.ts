// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRangeWithOptionalMax } from './number_range';


export interface IRecipeStepVesselCreationRequestInput {
   quantity: NumberRangeWithOptionalMax;
 vesselPreposition: string;
 name: string;
 notes: string;
 productOfRecipeStepIndex?: number;
 productOfRecipeStepProductIndex?: number;
 recipeStepProductID?: string;
 unavailableAfterStep: boolean;
 vesselID?: string;

}

export class RecipeStepVesselCreationRequestInput implements IRecipeStepVesselCreationRequestInput {
   quantity: NumberRangeWithOptionalMax;
 vesselPreposition: string;
 name: string;
 notes: string;
 productOfRecipeStepIndex?: number;
 productOfRecipeStepProductIndex?: number;
 recipeStepProductID?: string;
 unavailableAfterStep: boolean;
 vesselID?: string;
constructor(input: Partial<RecipeStepVesselCreationRequestInput> = {}) {
	 this.quantity = input.quantity = { min: 0 };
 this.vesselPreposition = input.vesselPreposition = '';
 this.name = input.name = '';
 this.notes = input.notes = '';
 this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
 this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
 this.recipeStepProductID = input.recipeStepProductID;
 this.unavailableAfterStep = input.unavailableAfterStep = false;
 this.vesselID = input.vesselID;
}
}