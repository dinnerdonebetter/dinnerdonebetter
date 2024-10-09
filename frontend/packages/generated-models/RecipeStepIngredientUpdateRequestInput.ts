// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepIngredientUpdateRequestInput {
  productOfRecipeID?: string;
  productPercentageToUse?: number;
  quantity: OptionalNumberRange;
  quantityNotes?: string;
  recipeStepProductID?: string;
  toTaste?: boolean;
  ingredientID?: string;
  name?: string;
  optional?: boolean;
  vesselIndex?: number;
  belongsToRecipeStep?: string;
  optionIndex?: number;
  ingredientNotes?: string;
  measurementUnitID?: string;
}

export class RecipeStepIngredientUpdateRequestInput implements IRecipeStepIngredientUpdateRequestInput {
  productOfRecipeID?: string;
  productPercentageToUse?: number;
  quantity: OptionalNumberRange;
  quantityNotes?: string;
  recipeStepProductID?: string;
  toTaste?: boolean;
  ingredientID?: string;
  name?: string;
  optional?: boolean;
  vesselIndex?: number;
  belongsToRecipeStep?: string;
  optionIndex?: number;
  ingredientNotes?: string;
  measurementUnitID?: string;
  constructor(input: Partial<RecipeStepIngredientUpdateRequestInput> = {}) {
    this.productOfRecipeID = input.productOfRecipeID;
    this.productPercentageToUse = input.productPercentageToUse;
    this.quantity = input.quantity = {};
    this.quantityNotes = input.quantityNotes;
    this.recipeStepProductID = input.recipeStepProductID;
    this.toTaste = input.toTaste;
    this.ingredientID = input.ingredientID;
    this.name = input.name;
    this.optional = input.optional;
    this.vesselIndex = input.vesselIndex;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
    this.optionIndex = input.optionIndex;
    this.ingredientNotes = input.ingredientNotes;
    this.measurementUnitID = input.measurementUnitID;
  }
}
