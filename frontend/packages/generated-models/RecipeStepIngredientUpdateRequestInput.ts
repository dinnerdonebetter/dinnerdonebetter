// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IRecipeStepIngredientUpdateRequestInput {
  ingredientNotes?: string;
  measurementUnitID?: string;
  name?: string;
  ingredientID?: string;
  productPercentageToUse?: number;
  recipeStepProductID?: string;
  optional?: boolean;
  quantityNotes?: string;
  quantity: OptionalNumberRange;
  optionIndex?: number;
  productOfRecipeID?: string;
  toTaste?: boolean;
  vesselIndex?: number;
  belongsToRecipeStep?: string;
}

export class RecipeStepIngredientUpdateRequestInput implements IRecipeStepIngredientUpdateRequestInput {
  ingredientNotes?: string;
  measurementUnitID?: string;
  name?: string;
  ingredientID?: string;
  productPercentageToUse?: number;
  recipeStepProductID?: string;
  optional?: boolean;
  quantityNotes?: string;
  quantity: OptionalNumberRange;
  optionIndex?: number;
  productOfRecipeID?: string;
  toTaste?: boolean;
  vesselIndex?: number;
  belongsToRecipeStep?: string;
  constructor(input: Partial<RecipeStepIngredientUpdateRequestInput> = {}) {
    this.ingredientNotes = input.ingredientNotes;
    this.measurementUnitID = input.measurementUnitID;
    this.name = input.name;
    this.ingredientID = input.ingredientID;
    this.productPercentageToUse = input.productPercentageToUse;
    this.recipeStepProductID = input.recipeStepProductID;
    this.optional = input.optional;
    this.quantityNotes = input.quantityNotes;
    this.quantity = input.quantity = {};
    this.optionIndex = input.optionIndex;
    this.productOfRecipeID = input.productOfRecipeID;
    this.toTaste = input.toTaste;
    this.vesselIndex = input.vesselIndex;
    this.belongsToRecipeStep = input.belongsToRecipeStep;
  }
}
