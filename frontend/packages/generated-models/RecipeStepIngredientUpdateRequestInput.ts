// GENERATED CODE, DO NOT EDIT MANUALLY

 import { OptionalNumberRange } from './number_range';


export interface IRecipeStepIngredientUpdateRequestInput {
   ingredientID?: string;
 optionIndex?: number;
 ingredientNotes?: string;
 productOfRecipeID?: string;
 productPercentageToUse?: number;
 quantity: OptionalNumberRange;
 toTaste?: boolean;
 vesselIndex?: number;
 belongsToRecipeStep?: string;
 measurementUnitID?: string;
 name?: string;
 optional?: boolean;
 quantityNotes?: string;
 recipeStepProductID?: string;

}

export class RecipeStepIngredientUpdateRequestInput implements IRecipeStepIngredientUpdateRequestInput {
   ingredientID?: string;
 optionIndex?: number;
 ingredientNotes?: string;
 productOfRecipeID?: string;
 productPercentageToUse?: number;
 quantity: OptionalNumberRange;
 toTaste?: boolean;
 vesselIndex?: number;
 belongsToRecipeStep?: string;
 measurementUnitID?: string;
 name?: string;
 optional?: boolean;
 quantityNotes?: string;
 recipeStepProductID?: string;
constructor(input: Partial<RecipeStepIngredientUpdateRequestInput> = {}) {
	 this.ingredientID = input.ingredientID;
 this.optionIndex = input.optionIndex;
 this.ingredientNotes = input.ingredientNotes;
 this.productOfRecipeID = input.productOfRecipeID;
 this.productPercentageToUse = input.productPercentageToUse;
 this.quantity = input.quantity = {};
 this.toTaste = input.toTaste;
 this.vesselIndex = input.vesselIndex;
 this.belongsToRecipeStep = input.belongsToRecipeStep;
 this.measurementUnitID = input.measurementUnitID;
 this.name = input.name;
 this.optional = input.optional;
 this.quantityNotes = input.quantityNotes;
 this.recipeStepProductID = input.recipeStepProductID;
}
}