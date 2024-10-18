// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { NumberRangeWithOptionalMax } from './number_range.gen';

export interface IRecipeStepIngredient {
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  ingredientNotes: string;
  lastUpdatedAt: string;
  measurementUnit: ValidMeasurementUnit;
  name: string;
  optionIndex: number;
  optional: boolean;
  productOfRecipeID: string;
  productPercentageToUse: number;
  quantity: NumberRangeWithOptionalMax;
  quantityNotes: string;
  recipeStepProductID: string;
  toTaste: boolean;
  vesselIndex: number;
}

export class RecipeStepIngredient implements IRecipeStepIngredient {
  archivedAt: string;
  belongsToRecipeStep: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  ingredientNotes: string;
  lastUpdatedAt: string;
  measurementUnit: ValidMeasurementUnit;
  name: string;
  optionIndex: number;
  optional: boolean;
  productOfRecipeID: string;
  productPercentageToUse: number;
  quantity: NumberRangeWithOptionalMax;
  quantityNotes: string;
  recipeStepProductID: string;
  toTaste: boolean;
  vesselIndex: number;
  constructor(input: Partial<RecipeStepIngredient> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToRecipeStep = input.belongsToRecipeStep || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.ingredientNotes = input.ingredientNotes || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.measurementUnit = input.measurementUnit || new ValidMeasurementUnit();
    this.name = input.name || '';
    this.optionIndex = input.optionIndex || 0;
    this.optional = input.optional || false;
    this.productOfRecipeID = input.productOfRecipeID || '';
    this.productPercentageToUse = input.productPercentageToUse || 0;
    this.quantity = input.quantity || { min: 0 };
    this.quantityNotes = input.quantityNotes || '';
    this.recipeStepProductID = input.recipeStepProductID || '';
    this.toTaste = input.toTaste || false;
    this.vesselIndex = input.vesselIndex || 0;
  }
}
