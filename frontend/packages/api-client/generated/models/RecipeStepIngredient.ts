/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ValidIngredient } from './ValidIngredient';
import type { ValidMeasurementUnit } from './ValidMeasurementUnit';
export type RecipeStepIngredient = {
  archivedAt?: string;
  belongsToRecipeStep?: string;
  createdAt?: string;
  id?: string;
  ingredient?: ValidIngredient;
  ingredientNotes?: string;
  lastUpdatedAt?: string;
  maximumQuantity?: number;
  measurementUnit?: ValidMeasurementUnit;
  minimumQuantity?: number;
  name?: string;
  optionIndex?: number;
  optional?: boolean;
  productOfRecipeID?: string;
  productPercentageToUse?: number;
  quantityNotes?: string;
  recipeStepProductID?: string;
  toTaste?: boolean;
  vesselIndex?: number;
};
