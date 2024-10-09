// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidMealPlanGroceryListItemStatus } from './enums';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealPlanGroceryListItem {
  purchasedUPC: string;
  quantityPurchased: number;
  archivedAt: string;
  createdAt: string;
  id: string;
  measurementUnit: ValidMeasurementUnit;
  purchasedMeasurementUnit: ValidMeasurementUnit;
  quantityNeeded: NumberRangeWithOptionalMax;
  belongsToMealPlan: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  purchasePrice: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
}

export class MealPlanGroceryListItem implements IMealPlanGroceryListItem {
  purchasedUPC: string;
  quantityPurchased: number;
  archivedAt: string;
  createdAt: string;
  id: string;
  measurementUnit: ValidMeasurementUnit;
  purchasedMeasurementUnit: ValidMeasurementUnit;
  quantityNeeded: NumberRangeWithOptionalMax;
  belongsToMealPlan: string;
  ingredient: ValidIngredient;
  lastUpdatedAt: string;
  purchasePrice: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  constructor(input: Partial<MealPlanGroceryListItem> = {}) {
    this.purchasedUPC = input.purchasedUPC || '';
    this.quantityPurchased = input.quantityPurchased || 0;
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.measurementUnit = input.measurementUnit || new ValidMeasurementUnit();
    this.purchasedMeasurementUnit = input.purchasedMeasurementUnit || new ValidMeasurementUnit();
    this.quantityNeeded = input.quantityNeeded || { min: 0 };
    this.belongsToMealPlan = input.belongsToMealPlan || '';
    this.ingredient = input.ingredient || new ValidIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.purchasePrice = input.purchasePrice || 0;
    this.status = input.status || 'unknown';
    this.statusExplanation = input.statusExplanation || '';
  }
}
