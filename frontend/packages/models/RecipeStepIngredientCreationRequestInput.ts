// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepIngredientCreationRequestInput {
  vesselIndex: number;
  measurementUnitID: string;
  optionIndex: number;
  quantity: NumberRangeWithOptionalMax;
  optional: boolean;
  quantityNotes: string;
  toTaste: boolean;
  productOfRecipeStepProductIndex: number;
  productPercentageToUse: number;
  ingredientID: string;
  name: string;
  productOfRecipeID: string;
  ingredientNotes: string;
  productOfRecipeStepIndex: number;
}

export class RecipeStepIngredientCreationRequestInput implements IRecipeStepIngredientCreationRequestInput {
  vesselIndex: number;
  measurementUnitID: string;
  optionIndex: number;
  quantity: NumberRangeWithOptionalMax;
  optional: boolean;
  quantityNotes: string;
  toTaste: boolean;
  productOfRecipeStepProductIndex: number;
  productPercentageToUse: number;
  ingredientID: string;
  name: string;
  productOfRecipeID: string;
  ingredientNotes: string;
  productOfRecipeStepIndex: number;
  constructor(input: Partial<RecipeStepIngredientCreationRequestInput> = {}) {
    this.vesselIndex = input.vesselIndex || 0;
    this.measurementUnitID = input.measurementUnitID || '';
    this.optionIndex = input.optionIndex || 0;
    this.quantity = input.quantity || { min: 0 };
    this.optional = input.optional || false;
    this.quantityNotes = input.quantityNotes || '';
    this.toTaste = input.toTaste || false;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex || 0;
    this.productPercentageToUse = input.productPercentageToUse || 0;
    this.ingredientID = input.ingredientID || '';
    this.name = input.name || '';
    this.productOfRecipeID = input.productOfRecipeID || '';
    this.ingredientNotes = input.ingredientNotes || '';
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex || 0;
  }
}
