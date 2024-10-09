// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepIngredientCreationRequestInput {
  productOfRecipeID?: string;
  quantity: NumberRangeWithOptionalMax;
  vesselIndex?: number;
  measurementUnitID: string;
  optional: boolean;
  quantityNotes: string;
  toTaste: boolean;
  ingredientID?: string;
  name: string;
  optionIndex: number;
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  ingredientNotes: string;
  productPercentageToUse?: number;
}

export class RecipeStepIngredientCreationRequestInput implements IRecipeStepIngredientCreationRequestInput {
  productOfRecipeID?: string;
  quantity: NumberRangeWithOptionalMax;
  vesselIndex?: number;
  measurementUnitID: string;
  optional: boolean;
  quantityNotes: string;
  toTaste: boolean;
  ingredientID?: string;
  name: string;
  optionIndex: number;
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  ingredientNotes: string;
  productPercentageToUse?: number;
  constructor(input: Partial<RecipeStepIngredientCreationRequestInput> = {}) {
    this.productOfRecipeID = input.productOfRecipeID;
    this.quantity = input.quantity = { min: 0 };
    this.vesselIndex = input.vesselIndex;
    this.measurementUnitID = input.measurementUnitID = '';
    this.optional = input.optional = false;
    this.quantityNotes = input.quantityNotes = '';
    this.toTaste = input.toTaste = false;
    this.ingredientID = input.ingredientID;
    this.name = input.name = '';
    this.optionIndex = input.optionIndex = 0;
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
    this.ingredientNotes = input.ingredientNotes = '';
    this.productPercentageToUse = input.productPercentageToUse;
  }
}
