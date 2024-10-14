// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepIngredientUpdateRequestInput {
  belongsToRecipeStep: string;
  ingredientID: string;
  ingredientNotes: string;
  measurementUnitID: string;
  name: string;
  optionIndex: number;
  optional: boolean;
  productOfRecipeID: string;
  productPercentageToUse: number;
  quantity: OptionalNumberRange;
  quantityNotes: string;
  recipeStepProductID: string;
  toTaste: boolean;
  vesselIndex: number;
}

export class RecipeStepIngredientUpdateRequestInput implements IRecipeStepIngredientUpdateRequestInput {
  belongsToRecipeStep: string;
  ingredientID: string;
  ingredientNotes: string;
  measurementUnitID: string;
  name: string;
  optionIndex: number;
  optional: boolean;
  productOfRecipeID: string;
  productPercentageToUse: number;
  quantity: OptionalNumberRange;
  quantityNotes: string;
  recipeStepProductID: string;
  toTaste: boolean;
  vesselIndex: number;
  constructor(input: Partial<RecipeStepIngredientUpdateRequestInput> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.ingredientID = input.ingredientID || '';
    this.ingredientNotes = input.ingredientNotes || '';
    this.measurementUnitID = input.measurementUnitID || '';
    this.name = input.name || '';
    this.optionIndex = input.optionIndex || 0;
    this.optional = input.optional || false;
    this.productOfRecipeID = input.productOfRecipeID || '';
    this.productPercentageToUse = input.productPercentageToUse || 0;
    this.quantity = input.quantity || {};
    this.quantityNotes = input.quantityNotes || '';
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.toTaste = input.toTaste || false;
    this.vesselIndex = input.vesselIndex || 0;
  }
}
