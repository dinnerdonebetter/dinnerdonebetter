// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient';
import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidMealPlanGroceryListItemStatus } from './enums';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealPlanGroceryListItem {
  statusExplanation: string;
  id: string;
  measurementUnit: ValidMeasurementUnit;
  archivedAt?: string;
  createdAt: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased?: number;
  belongsToMealPlan: string;
  lastUpdatedAt?: string;
  purchasePrice?: number;
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  purchasedUPC?: string;
  status: ValidMealPlanGroceryListItemStatus;
  ingredient: ValidIngredient;
}

export class MealPlanGroceryListItem implements IMealPlanGroceryListItem {
  statusExplanation: string;
  id: string;
  measurementUnit: ValidMeasurementUnit;
  archivedAt?: string;
  createdAt: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased?: number;
  belongsToMealPlan: string;
  lastUpdatedAt?: string;
  purchasePrice?: number;
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  purchasedUPC?: string;
  status: ValidMealPlanGroceryListItemStatus;
  ingredient: ValidIngredient;
  constructor(input: Partial<MealPlanGroceryListItem> = {}) {
    this.statusExplanation = input.statusExplanation = '';
    this.id = input.id = '';
    this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.quantityNeeded = input.quantityNeeded = { min: 0 };
    this.quantityPurchased = input.quantityPurchased;
    this.belongsToMealPlan = input.belongsToMealPlan = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.purchasePrice = input.purchasePrice;
    this.purchasedMeasurementUnit = input.purchasedMeasurementUnit;
    this.purchasedUPC = input.purchasedUPC;
    this.status = input.status = 'unknown';
    this.ingredient = input.ingredient = new ValidIngredient();
  }
}
