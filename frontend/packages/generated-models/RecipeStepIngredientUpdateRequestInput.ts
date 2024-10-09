// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepIngredientUpdateRequestInput {
  name?: string;
  productOfRecipeID?: string;
  recipeStepProductID?: string;
  vesselIndex?: number;
  measurementUnitID?: string;
  optional?: boolean;
  optionIndex?: number;
  quantityNotes?: string;
  productPercentageToUse?: number;
  ingredientID?: string;
  ingredientNotes?: string;
  quantity: OptionalNumberRange;
  toTaste?: boolean;
  belongsToRecipeStep?: string;
}

export class RecipeStepIngredientUpdateRequestInput implements IRecipeStepIngredientUpdateRequestInput {
  name?: string;
  productOfRecipeID?: string;
  recipeStepProductID?: string;
  vesselIndex?: number;
  measurementUnitID?: string;
  optional?: boolean;
  optionIndex?: number;
  quantityNotes?: string;
  productPercentageToUse?: number;
  ingredientID?: string;
  ingredientNotes?: string;
  quantity: OptionalNumberRange;
  toTaste?: boolean;
  belongsToRecipeStep?: string;
  constructor(input: Partial<RecipeStepIngredientUpdateRequestInput> = {}) {
    this.name = input.name;
    this.productOfRecipeID = input.productOfRecipeID;
    this.recipeStepProductID = input.recipeStepProductID;
    this.vesselIndex = input.vesselIndex;
    this.measurementUnitID = input.measurementUnitID;
    this.optional = input.optional;
    this.optionIndex = input.optionIndex;
    this.quantityNotes = input.quantityNotes;
    this.productPercentageToUse = input.productPercentageToUse;
    this.ingredientID = input.ingredientID;
    this.ingredientNotes = input.ingredientNotes;
    this.quantity = input.quantity = {};
    this.toTaste = input.toTaste;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
  }
}
