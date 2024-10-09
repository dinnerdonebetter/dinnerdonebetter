// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepIngredient {
  ingredient?: ValidIngredient;
  ingredientNotes: string;
  productOfRecipeID?: string;
  productPercentageToUse?: number;
  quantityNotes: string;
  id: string;
  measurementUnit: ValidMeasurementUnit;
  optional: boolean;
  archivedAt?: string;
  belongsToRecipeStep: string;
  name: string;
  optionIndex: number;
  recipeStepProductID?: string;
  createdAt: string;
  lastUpdatedAt?: string;
  quantity: NumberRangeWithOptionalMax;
  toTaste: boolean;
  vesselIndex?: number;
}

export class RecipeStepIngredient implements IRecipeStepIngredient {
  ingredient?: ValidIngredient;
  ingredientNotes: string;
  productOfRecipeID?: string;
  productPercentageToUse?: number;
  quantityNotes: string;
  id: string;
  measurementUnit: ValidMeasurementUnit;
  optional: boolean;
  archivedAt?: string;
  belongsToRecipeStep: string;
  name: string;
  optionIndex: number;
  recipeStepProductID?: string;
  createdAt: string;
  lastUpdatedAt?: string;
  quantity: NumberRangeWithOptionalMax;
  toTaste: boolean;
  vesselIndex?: number;
  constructor(input: Partial<RecipeStepIngredient> = {}) {
    this.ingredient = input.ingredient;
    this.ingredientNotes = input.ingredientNotes = '';
    this.productOfRecipeID = input.productOfRecipeID;
    this.productPercentageToUse = input.productPercentageToUse;
    this.quantityNotes = input.quantityNotes = '';
    this.id = input.id = '';
    this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
    this.optional = input.optional = false;
    this.archivedAt = input.archivedAt;
    this.belongsToRecipeStep = input.belongsToRecipeStep = '';
    this.name = input.name = '';
    this.optionIndex = input.optionIndex = 0;
    this.recipeStepProductID = input.recipeStepProductID;
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.quantity = input.quantity = { min: 0 };
    this.toTaste = input.toTaste = false;
    this.vesselIndex = input.vesselIndex;
  }
}
