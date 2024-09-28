/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ValidIngredient } from './ValidIngredient';
import type { ValidMeasurementUnit } from './ValidMeasurementUnit';
export type MealPlanGroceryListItem = {
  archivedAt?: string;
  belongsToMealPlan?: string;
  createdAt?: string;
  id?: string;
  ingredient?: ValidIngredient;
  lastUpdatedAt?: string;
  maximumQuantityNeeded?: number;
  measurementUnit?: ValidMeasurementUnit;
  minimumQuantityNeeded?: number;
  purchasePrice?: number;
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  purchasedUPC?: string;
  quantityPurchased?: number;
  status?: string;
  statusExplanation?: string;
};
