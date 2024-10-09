// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepIngredientCreationRequestInput {
  ingredientID?: string;
  productPercentageToUse?: number;
  measurementUnitID: string;
  productOfRecipeID?: string;
  productOfRecipeStepIndex?: number;
  vesselIndex?: number;
  ingredientNotes: string;
  optional: boolean;
  quantity: NumberRangeWithOptionalMax;
  name: string;
  optionIndex: number;
  productOfRecipeStepProductIndex?: number;
  quantityNotes: string;
  toTaste: boolean;
}

export class RecipeStepIngredientCreationRequestInput implements IRecipeStepIngredientCreationRequestInput {
  ingredientID?: string;
  productPercentageToUse?: number;
  measurementUnitID: string;
  productOfRecipeID?: string;
  productOfRecipeStepIndex?: number;
  vesselIndex?: number;
  ingredientNotes: string;
  optional: boolean;
  quantity: NumberRangeWithOptionalMax;
  name: string;
  optionIndex: number;
  productOfRecipeStepProductIndex?: number;
  quantityNotes: string;
  toTaste: boolean;
  constructor(input: Partial<RecipeStepIngredientCreationRequestInput> = {}) {
    this.ingredientID = input.ingredientID;
    this.productPercentageToUse = input.productPercentageToUse;
    this.measurementUnitID = input.measurementUnitID = '';
    this.productOfRecipeID = input.productOfRecipeID;
    this.productOfRecipeStepIndex = input.productOfRecipeStepIndex;
    this.vesselIndex = input.vesselIndex;
    this.ingredientNotes = input.ingredientNotes = '';
    this.optional = input.optional = false;
    this.quantity = input.quantity = { min: 0 };
    this.name = input.name = '';
    this.optionIndex = input.optionIndex = 0;
    this.productOfRecipeStepProductIndex = input.productOfRecipeStepProductIndex;
    this.quantityNotes = input.quantityNotes = '';
    this.toTaste = input.toTaste = false;
  }
}
