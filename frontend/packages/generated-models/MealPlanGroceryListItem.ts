// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidMealPlanGroceryListItemStatus } from './enums';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealPlanGroceryListItem {
  measurementUnit: ValidMeasurementUnit;
  belongsToMealPlan: string;
  createdAt: string;
  status: ValidMealPlanGroceryListItemStatus;
  id: string;
  purchasePrice?: number;
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  statusExplanation: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased?: number;
  archivedAt?: string;
  ingredient: ValidIngredient;
  lastUpdatedAt?: string;
  purchasedUPC?: string;
}

export class MealPlanGroceryListItem implements IMealPlanGroceryListItem {
  measurementUnit: ValidMeasurementUnit;
  belongsToMealPlan: string;
  createdAt: string;
  status: ValidMealPlanGroceryListItemStatus;
  id: string;
  purchasePrice?: number;
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  statusExplanation: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased?: number;
  archivedAt?: string;
  ingredient: ValidIngredient;
  lastUpdatedAt?: string;
  purchasedUPC?: string;
  constructor(input: Partial<MealPlanGroceryListItem> = {}) {
    this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
    this.belongsToMealPlan = input.belongsToMealPlan = '';
    this.createdAt = input.createdAt = '';
    this.status = input.status = 'unknown';
    this.id = input.id = '';
    this.purchasePrice = input.purchasePrice;
    this.purchasedMeasurementUnit = input.purchasedMeasurementUnit;
    this.statusExplanation = input.statusExplanation = '';
    this.quantityNeeded = input.quantityNeeded = { min: 0 };
    this.quantityPurchased = input.quantityPurchased;
    this.archivedAt = input.archivedAt;
    this.ingredient = input.ingredient = new ValidIngredient();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.purchasedUPC = input.purchasedUPC;
  }
}
