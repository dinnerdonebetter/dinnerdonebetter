// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepIngredientUpdateRequestInput {
  productPercentageToUse: number;
  quantity: OptionalNumberRange;
  toTaste: boolean;
  vesselIndex: number;
  ingredientID: string;
  ingredientNotes: string;
  productOfRecipeID: string;
  name: string;
  optional: boolean;
  quantityNotes: string;
  belongsToRecipeStep: string;
  measurementUnitID: string;
  optionIndex: number;
  recipeStepProductID: string;
}

export class RecipeStepIngredientUpdateRequestInput implements IRecipeStepIngredientUpdateRequestInput {
  productPercentageToUse: number;
  quantity: OptionalNumberRange;
  toTaste: boolean;
  vesselIndex: number;
  ingredientID: string;
  ingredientNotes: string;
  productOfRecipeID: string;
  name: string;
  optional: boolean;
  quantityNotes: string;
  belongsToRecipeStep: string;
  measurementUnitID: string;
  optionIndex: number;
  recipeStepProductID: string;
  constructor(input: Partial<RecipeStepIngredientUpdateRequestInput> = {}) {
    this.productPercentageToUse = input.productPercentageToUse || 0;
    this.quantity = input.quantity || {};
    this.toTaste = input.toTaste || false;
    this.vesselIndex = input.vesselIndex || 0;
    this.ingredientID = input.ingredientID || '';
    this.ingredientNotes = input.ingredientNotes || '';
    this.productOfRecipeID = input.productOfRecipeID || '';
    this.name = input.name || '';
    this.optional = input.optional || false;
    this.quantityNotes = input.quantityNotes || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.measurementUnitID = input.measurementUnitID || '';
    this.optionIndex = input.optionIndex || 0;
    this.recipeStepProductID = input.recipeStepProductID || '';
  }
}
