// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient.gen';
import { ValidMeasurementUnit } from './ValidMeasurementUnit.gen';
import { ValidMealPlanGroceryListItemStatus } from './enums.gen';
import { NumberRangeWithOptionalMax } from './number_range.gen';

export interface IMealPlanGroceryListItem {
  archivedAt: string;
  belongsToMealPlan: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  measurementUnit: ValidMeasurementUnit;
  purchasePrice: number;
  purchasedMeasurementUnit: ValidMeasurementUnit;
  purchasedUPC: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
}

export class MealPlanGroceryListItem implements IMealPlanGroceryListItem {
  archivedAt: string;
  belongsToMealPlan: string;
  createdAt: string;
  id: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  measurementUnit: ValidMeasurementUnit;
  purchasePrice: number;
  purchasedMeasurementUnit: ValidMeasurementUnit;
  purchasedUPC: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  constructor(input: Partial<MealPlanGroceryListItem> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToMealPlan = input.belongsToMealPlan || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.measurementUnit = input.measurementUnit || new ValidMeasurementUnit();
    this.purchasePrice = input.purchasePrice || 0;
    this.purchasedMeasurementUnit = input.purchasedMeasurementUnit || new ValidMeasurementUnit();
    this.purchasedUPC = input.purchasedUPC || '';
    this.quantityNeeded = input.quantityNeeded || { min: 0 };
    this.quantityPurchased = input.quantityPurchased || 0;
    this.status = input.status || 'unknown';
    this.statusExplanation = input.statusExplanation || '';
  }
}
