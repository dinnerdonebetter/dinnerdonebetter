// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepIngredient {
  archivedAt?: string;
  productPercentageToUse?: number;
  quantityNotes: string;
  belongsToRecipeStep: string;
  ingredient?: ValidIngredient;
  measurementUnit: ValidMeasurementUnit;
  name: string;
  quantity: NumberRangeWithOptionalMax;
  toTaste: boolean;
  vesselIndex?: number;
  createdAt: string;
  lastUpdatedAt?: string;
  optionIndex: number;
  productOfRecipeID?: string;
  id: string;
  ingredientNotes: string;
  optional: boolean;
  recipeStepProductID?: string;
}

export class RecipeStepIngredient implements IRecipeStepIngredient {
  archivedAt?: string;
  productPercentageToUse?: number;
  quantityNotes: string;
  belongsToRecipeStep: string;
  ingredient?: ValidIngredient;
  measurementUnit: ValidMeasurementUnit;
  name: string;
  quantity: NumberRangeWithOptionalMax;
  toTaste: boolean;
  vesselIndex?: number;
  createdAt: string;
  lastUpdatedAt?: string;
  optionIndex: number;
  productOfRecipeID?: string;
  id: string;
  ingredientNotes: string;
  optional: boolean;
  recipeStepProductID?: string;
  constructor(input: Partial<RecipeStepIngredient> = {}) {
    this.archivedAt = input.archivedAt;
    this.productPercentageToUse = input.productPercentageToUse;
    this.quantityNotes = input.quantityNotes = '';
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.ingredient = input.ingredient;
    this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
    this.name = input.name = '';
    this.quantity = input.quantity = { min: 0 };
    this.toTaste = input.toTaste = false;
    this.vesselIndex = input.vesselIndex;
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.optionIndex = input.optionIndex = 0;
    this.productOfRecipeID = input.productOfRecipeID;
    this.id = input.id = '';
    this.ingredientNotes = input.ingredientNotes = '';
    this.optional = input.optional = false;
    this.recipeStepProductID = input.recipeStepProductID;
  }
}
