// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IRecipeStepIngredient {
  belongsToRecipeStep: string;
  optionIndex: number;
  optional: boolean;
  name: string;
  productOfRecipeID: string;
  recipeStepProductID: string;
  toTaste: boolean;
  archivedAt: string;
  id: string;
  ingredientNotes: string;
  lastUpdatedAt: string;
  quantity: NumberRangeWithOptionalMax;
  quantityNotes: string;
  vesselIndex: number;
  createdAt: string;
  ingredient: ValidIngredient;
  measurementUnit: ValidMeasurementUnit;
  productPercentageToUse: number;
}

export class RecipeStepIngredient implements IRecipeStepIngredient {
  belongsToRecipeStep: string;
  optionIndex: number;
  optional: boolean;
  name: string;
  productOfRecipeID: string;
  recipeStepProductID: string;
  toTaste: boolean;
  archivedAt: string;
  id: string;
  ingredientNotes: string;
  lastUpdatedAt: string;
  quantity: NumberRangeWithOptionalMax;
  quantityNotes: string;
  vesselIndex: number;
  createdAt: string;
  ingredient: ValidIngredient;
  measurementUnit: ValidMeasurementUnit;
  productPercentageToUse: number;
  constructor(input: Partial<RecipeStepIngredient> = {}) {
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.optionIndex = input.optionIndex || 0;
    this.optional = input.optional || false;
    this.name = input.name || '';
    this.productOfRecipeID = input.productOfRecipeID || '';
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.toTaste = input.toTaste || false;
    this.archivedAt = input.archivedAt || '';
    this.id = input.id || '';
    this.ingredientNotes = input.ingredientNotes || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.quantity = input.quantity || { min: 0 };
    this.quantityNotes = input.quantityNotes || '';
    this.vesselIndex = input.vesselIndex || 0;
    this.createdAt = input.createdAt || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.measurementUnit = input.measurementUnit || new ValidMeasurementUnit();
    this.productPercentageToUse = input.productPercentageToUse || 0;
  }
}
