// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidMealPlanGroceryListItemStatus } from './enums';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealPlanGroceryListItem {
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  archivedAt?: string;
  belongsToMealPlan: string;
  createdAt: string;
  measurementUnit: ValidMeasurementUnit;
  ingredient: ValidIngredient;
  lastUpdatedAt?: string;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  id: string;
  purchasePrice?: number;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased?: number;
  purchasedUPC?: string;
}

export class MealPlanGroceryListItem implements IMealPlanGroceryListItem {
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  archivedAt?: string;
  belongsToMealPlan: string;
  createdAt: string;
  measurementUnit: ValidMeasurementUnit;
  ingredient: ValidIngredient;
  lastUpdatedAt?: string;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  id: string;
  purchasePrice?: number;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased?: number;
  purchasedUPC?: string;
  constructor(input: Partial<MealPlanGroceryListItem> = {}) {
    this.purchasedMeasurementUnit = input.purchasedMeasurementUnit;
    this.archivedAt = input.archivedAt;
    this.belongsToMealPlan = input.belongsToMealPlan = '';
    this.createdAt = input.createdAt = '';
    this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
    this.ingredient = input.ingredient = new ValidIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.status = input.status = 'unknown';
    this.statusExplanation = input.statusExplanation = '';
    this.id = input.id = '';
    this.purchasePrice = input.purchasePrice;
    this.quantityNeeded = input.quantityNeeded = { min: 0 };
    this.quantityPurchased = input.quantityPurchased;
    this.purchasedUPC = input.purchasedUPC;
  }
}
