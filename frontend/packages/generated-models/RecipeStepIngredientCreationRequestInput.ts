// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRangeWithOptionalMax } from './number_range';


export interface IRecipeStepIngredientCreationRequestInput {
   quantityNotes: string;
 ingredientID?: string;
 productOfRecipeID?: string;
 optional: boolean;
 vesselIndex?: number;
 productOfRecipeStepIndex?: number;
 productOfRecipeStepProductIndex?: number;
 toTaste: boolean;
 name: string;
 optionIndex: number;
 productPercentageToUse?: number;
 quantity: NumberRangeWithOptionalMax;
 ingredientNotes: string;
 measurementUnitID: string;

}

export class RecipeStepIngredientCreationRequestInput implements IRecipeStepIngredientCreationRequestInput {
   quantityNotes: string;
 ingredientID?: string;
 productOfRecipeID?: string;
 optional: boolean;
 vesselIndex?: number;
 productOfRecipeStepIndex?: number;
 productOfRecipeStepProductIndex?: number;
 toTaste: boolean;
 name: string;
 optionIndex: number;
 productPercentageToUse?: number;
 quantity: NumberRangeWithOptionalMax;
 ingredientNotes: string;
 measurementUnitID: string;
constructor(input: Partial<RecipeStepIngredientCreationRequestInput> = {}) {
	 this.quantityNotes = input.quantityNotes = '';
 this.ingredientID = input.ingredientID;
 this.productOfRecipeID = input.productOfRecipeID;
 this.optional = input.optional = false;
 this.vesselIndex = input.vesselIndex;
 this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
 this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
 this.toTaste = input.toTaste = false;
 this.name = input.name = '';
 this.optionIndex = input.optionIndex = 0;
 this.productPercentageToUse = input.productPercentageToUse;
 this.quantity = input.quantity = { min: 0 };
 this.ingredientNotes = input.ingredientNotes = '';
 this.measurementUnitID = input.measurementUnitID = '';
}
}