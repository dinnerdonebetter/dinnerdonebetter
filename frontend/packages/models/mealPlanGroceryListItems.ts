// Code generated by gen_typescript. DO NOT EDIT.

import { ValidMealPlanGroceryListItemStatus } from './_unions';
import { NumberRangeWithOptionalMax, OptionalNumberRange } from './main';
import { ValidIngredient } from './validIngredients';
import { ValidMeasurementUnit } from './validMeasurementUnits';

export interface IMealPlanGroceryListItem {
  createdAt: NonNullable<string>;
  quantityPurchased?: number;
  purchasePrice?: number;
  purchasedUPC?: string;
  archivedAt?: string;
  lastUpdatedAt?: string;
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  belongsToMealPlan: NonNullable<string>;
  status: NonNullable<ValidMealPlanGroceryListItemStatus>;
  statusExplanation: NonNullable<string>;
  id: NonNullable<string>;
  quantityNeeded: NonNullable<NumberRangeWithOptionalMax>;
  measurementUnit: NonNullable<ValidMeasurementUnit>;
  ingredient: NonNullable<ValidIngredient>;
}

export class MealPlanGroceryListItem implements IMealPlanGroceryListItem {
  createdAt: NonNullable<string> = '1970-01-01T00:00:00Z';
  quantityPurchased?: number;
  purchasePrice?: number;
  purchasedUPC?: string;
  archivedAt?: string;
  lastUpdatedAt?: string;
  purchasedMeasurementUnit?: ValidMeasurementUnit;
  belongsToMealPlan: NonNullable<string> = '';
  status: NonNullable<ValidMealPlanGroceryListItemStatus> = 'unknown';
  statusExplanation: NonNullable<string> = '';
  id: NonNullable<string> = '';
  quantityNeeded: NonNullable<NumberRangeWithOptionalMax> = { min: 0 };
  measurementUnit: NonNullable<ValidMeasurementUnit> = new ValidMeasurementUnit();
  ingredient: NonNullable<ValidIngredient> = new ValidIngredient();

  constructor(input: Partial<MealPlanGroceryListItem> = {}) {
    this.createdAt = input.createdAt ?? '1970-01-01T00:00:00Z';
    this.quantityPurchased = input.quantityPurchased;
    this.purchasePrice = input.purchasePrice;
    this.purchasedUPC = input.purchasedUPC;
    this.archivedAt = input.archivedAt;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.purchasedMeasurementUnit = input.purchasedMeasurementUnit;
    this.belongsToMealPlan = input.belongsToMealPlan ?? '';
    this.status = input.status ?? 'unknown';
    this.statusExplanation = input.statusExplanation ?? '';
    this.id = input.id ?? '';
    this.quantityNeeded = input.quantityNeeded ?? { min: 0 };
    this.measurementUnit = input.measurementUnit ?? new ValidMeasurementUnit();
    this.ingredient = input.ingredient ?? new ValidIngredient();
  }
}

export interface IMealPlanGroceryListItemCreationRequestInput {
  purchasedMeasurementUnitID?: string;
  purchasedUPC?: string;
  purchasePrice?: number;
  quantityPurchased?: number;
  status: NonNullable<ValidMealPlanGroceryListItemStatus>;
  belongsToMealPlan: NonNullable<string>;
  validIngredientID: NonNullable<string>;
  validMeasurementUnitID: NonNullable<string>;
  statusExplanation: NonNullable<string>;
  quantityNeeded: NonNullable<NumberRangeWithOptionalMax>;
}

export class MealPlanGroceryListItemCreationRequestInput implements IMealPlanGroceryListItemCreationRequestInput {
  purchasedMeasurementUnitID?: string;
  purchasedUPC?: string;
  purchasePrice?: number;
  quantityPurchased?: number;
  status: NonNullable<ValidMealPlanGroceryListItemStatus> = 'unknown';
  belongsToMealPlan: NonNullable<string> = '';
  validIngredientID: NonNullable<string> = '';
  validMeasurementUnitID: NonNullable<string> = '';
  statusExplanation: NonNullable<string> = '';
  quantityNeeded: NonNullable<NumberRangeWithOptionalMax> = { min: 0 };

  constructor(input: Partial<MealPlanGroceryListItemCreationRequestInput> = {}) {
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
    this.purchasedUPC = input.purchasedUPC;
    this.purchasePrice = input.purchasePrice;
    this.quantityPurchased = input.quantityPurchased;
    this.status = input.status ?? 'unknown';
    this.belongsToMealPlan = input.belongsToMealPlan ?? '';
    this.validIngredientID = input.validIngredientID ?? '';
    this.validMeasurementUnitID = input.validMeasurementUnitID ?? '';
    this.statusExplanation = input.statusExplanation ?? '';
    this.quantityNeeded = input.quantityNeeded ?? { min: 0 };
  }
}

export interface IMealPlanGroceryListItemUpdateRequestInput {
  belongsToMealPlan?: string;
  validIngredientID?: string;
  validMeasurementUnitID?: string;
  statusExplanation?: string;
  quantityPurchased?: number;
  purchasedMeasurementUnitID?: string;
  purchasedUPC?: string;
  purchasePrice?: number;
  status?: ValidMealPlanGroceryListItemStatus;
  quantityNeeded: NonNullable<OptionalNumberRange>;
}

export class MealPlanGroceryListItemUpdateRequestInput implements IMealPlanGroceryListItemUpdateRequestInput {
  belongsToMealPlan?: string;
  validIngredientID?: string;
  validMeasurementUnitID?: string;
  statusExplanation?: string;
  quantityPurchased?: number;
  purchasedMeasurementUnitID?: string;
  purchasedUPC?: string;
  purchasePrice?: number;
  status?: ValidMealPlanGroceryListItemStatus = 'unknown';
  quantityNeeded: NonNullable<OptionalNumberRange> = {};

  constructor(input: Partial<MealPlanGroceryListItemUpdateRequestInput> = {}) {
    this.belongsToMealPlan = input.belongsToMealPlan;
    this.validIngredientID = input.validIngredientID;
    this.validMeasurementUnitID = input.validMeasurementUnitID;
    this.statusExplanation = input.statusExplanation;
    this.quantityPurchased = input.quantityPurchased;
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
    this.purchasedUPC = input.purchasedUPC;
    this.purchasePrice = input.purchasePrice;
    this.status = input.status ?? 'unknown';
    this.quantityNeeded = input.quantityNeeded ?? {};
  }
}
