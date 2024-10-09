// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidIngredient } from './ValidIngredient';
 import { ValidMeasurementUnit } from './ValidMeasurementUnit';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IRecipeStepIngredient {
   name: string;
 quantity: NumberRangeWithOptionalMax;
 measurementUnit: ValidMeasurementUnit;
 optionIndex: number;
 optional: boolean;
 quantityNotes: string;
 recipeStepProductID?: string;
 ingredientNotes: string;
 ingredient?: ValidIngredient;
 lastUpdatedAt?: string;
 productPercentageToUse?: number;
 id: string;
 belongsToRecipeStep: string;
 createdAt: string;
 productOfRecipeID?: string;
 toTaste: boolean;
 vesselIndex?: number;
 archivedAt?: string;

}

export class RecipeStepIngredient implements IRecipeStepIngredient {
   name: string;
 quantity: NumberRangeWithOptionalMax;
 measurementUnit: ValidMeasurementUnit;
 optionIndex: number;
 optional: boolean;
 quantityNotes: string;
 recipeStepProductID?: string;
 ingredientNotes: string;
 ingredient?: ValidIngredient;
 lastUpdatedAt?: string;
 productPercentageToUse?: number;
 id: string;
 belongsToRecipeStep: string;
 createdAt: string;
 productOfRecipeID?: string;
 toTaste: boolean;
 vesselIndex?: number;
 archivedAt?: string;
constructor(input: Partial<RecipeStepIngredient> = {}) {
	 this.name = input.name = '';
 this.quantity = input.quantity = { min: 0 };
 this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
 this.optionIndex = input.optionIndex = 0;
 this.optional = input.optional = false;
 this.quantityNotes = input.quantityNotes = '';
 this.recipeStepProductID = input.recipeStepProductID;
 this.ingredientNotes = input.ingredientNotes = '';
 this.ingredient = input.ingredient;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.productPercentageToUse = input.productPercentageToUse;
 this.id = input.id = '';
 this.belongsToRecipeStep = input.belongsToRecipeStep = '';
 this.createdAt = input.createdAt = '';
 this.productOfRecipeID = input.productOfRecipeID;
 this.toTaste = input.toTaste = false;
 this.vesselIndex = input.vesselIndex;
 this.archivedAt = input.archivedAt;
}
}