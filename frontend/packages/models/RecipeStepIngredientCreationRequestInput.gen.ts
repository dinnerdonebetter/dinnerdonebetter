// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range.gen';

export interface IRecipeStepIngredientCreationRequestInput {
  ingredientID: string;
  ingredientNotes: string;
  measurementUnitID: string;
  name: string;
  optionIndex: number;
  optional: boolean;
  productOfRecipeID: string;
  productOfRecipeStepIndex: number;
  productOfRecipeStepProductIndex: number;
  productPercentageToUse: number;
  quantity: NumberRangeWithOptionalMax;
  quantityNotes: string;
  toTaste: boolean;
  vesselIndex: number;
}

export class RecipeStepIngredientCreationRequestInput implements IRecipeStepIngredientCreationRequestInput {
  ingredientID: string;
  ingredientNotes: string;
  measurementUnitID: string;
  name: string;
  optionIndex: number;
  optional: boolean;
  productOfRecipeID: string;
  productOfRecipeStepIndex: number;
  productOfRecipeStepProductIndex: number;
  productPercentageToUse: number;
  quantity: NumberRangeWithOptionalMax;
  quantityNotes: string;
  toTaste: boolean;
  vesselIndex: number;
  constructor(input: Partial<RecipeStepIngredientCreationRequestInput> = {}) {
    this.ingredientID = input.ingredientID || '';
    this.ingredientNotes = input.ingredientNotes || '';
    this.measurementUnitID = input.measurementUnitID || '';
    this.name = input.name || '';
    this.optionIndex = input.optionIndex || 0;
    this.optional = input.optional || false;
    this.productOfRecipeID = input.productOfRecipeID || '';
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex || 0;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex || 0;
    this.productPercentageToUse = input.productPercentageToUse || 0;
    this.quantity = input.quantity || { min: 0 };
    this.quantityNotes = input.quantityNotes || '';
    this.toTaste = input.toTaste || false;
    this.vesselIndex = input.vesselIndex || 0;
  }
}
