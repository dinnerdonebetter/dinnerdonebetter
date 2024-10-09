// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepIngredient {
  optionIndex: number;
  quantityNotes: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  measurementUnit: ValidMeasurementUnit;
  vesselIndex?: number;
  archivedAt?: string;
  optional: boolean;
  productPercentageToUse?: number;
  toTaste: boolean;
  lastUpdatedAt?: string;
  name: string;
  recipeStepProductID?: string;
  ingredient?: ValidIngredient;
  ingredientNotes: string;
  productOfRecipeID?: string;
  quantity: NumberRangeWithOptionalMax;
}

export class RecipeStepIngredient implements IRecipeStepIngredient {
  optionIndex: number;
  quantityNotes: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  measurementUnit: ValidMeasurementUnit;
  vesselIndex?: number;
  archivedAt?: string;
  optional: boolean;
  productPercentageToUse?: number;
  toTaste: boolean;
  lastUpdatedAt?: string;
  name: string;
  recipeStepProductID?: string;
  ingredient?: ValidIngredient;
  ingredientNotes: string;
  productOfRecipeID?: string;
  quantity: NumberRangeWithOptionalMax;
  constructor(input: Partial<RecipeStepIngredient> = {}) {
    this.optionIndex = input.optionIndex = 0;
    this.quantityNotes = input.quantityNotes = '';
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
    this.vesselIndex = input.vesselIndex;
    this.archivedAt = input.archivedAt;
    this.optional = input.optional = false;
    this.productPercentageToUse = input.productPercentageToUse;
    this.toTaste = input.toTaste = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.name = input.name = '';
    this.recipeStepProductID = input.recipeStepProductID;
    this.ingredient = input.ingredient;
    this.ingredientNotes = input.ingredientNotes = '';
    this.productOfRecipeID = input.productOfRecipeID;
    this.quantity = input.quantity = { min: 0 };
  }
}
