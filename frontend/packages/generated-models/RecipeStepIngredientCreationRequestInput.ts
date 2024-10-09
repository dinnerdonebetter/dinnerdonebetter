// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepIngredientCreationRequestInput {
  optional: boolean;
  quantity: NumberRangeWithOptionalMax;
  quantityNotes: string;
  name: string;
  optionIndex: number;
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  ingredientID?: string;
  productOfRecipeID?: string;
  toTaste: boolean;
  vesselIndex?: number;
  measurementUnitID: string;
  productPercentageToUse?: number;
  ingredientNotes: string;
}

export class RecipeStepIngredientCreationRequestInput implements IRecipeStepIngredientCreationRequestInput {
  optional: boolean;
  quantity: NumberRangeWithOptionalMax;
  quantityNotes: string;
  name: string;
  optionIndex: number;
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  ingredientID?: string;
  productOfRecipeID?: string;
  toTaste: boolean;
  vesselIndex?: number;
  measurementUnitID: string;
  productPercentageToUse?: number;
  ingredientNotes: string;
  constructor(input: Partial<RecipeStepIngredientCreationRequestInput> = {}) {
    this.optional = input.optional = false;
    this.quantity = input.quantity = { min: 0 };
    this.quantityNotes = input.quantityNotes = '';
    this.name = input.name = '';
    this.optionIndex = input.optionIndex = 0;
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
    this.ingredientID = input.ingredientID;
    this.productOfRecipeID = input.productOfRecipeID;
    this.toTaste = input.toTaste = false;
    this.vesselIndex = input.vesselIndex;
    this.measurementUnitID = input.measurementUnitID = '';
    this.productPercentageToUse = input.productPercentageToUse;
    this.ingredientNotes = input.ingredientNotes = '';
  }
}
