// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealPlanGroceryListItemCreationRequestInput {
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased?: number;
  validIngredientID: string;
  purchasedMeasurementUnitID?: string;
  purchasePrice?: number;
  purchasedUPC?: string;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  validMeasurementUnitID: string;
  belongsToMealPlan: string;
}

export class MealPlanGroceryListItemCreationRequestInput implements IMealPlanGroceryListItemCreationRequestInput {
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased?: number;
  validIngredientID: string;
  purchasedMeasurementUnitID?: string;
  purchasePrice?: number;
  purchasedUPC?: string;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  validMeasurementUnitID: string;
  belongsToMealPlan: string;
  constructor(input: Partial<MealPlanGroceryListItemCreationRequestInput> = {}) {
    this.quantityNeeded = input.quantityNeeded = { min: 0 };
    this.quantityPurchased = input.quantityPurchased;
    this.validIngredientID = input.validIngredientID = '';
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
    this.purchasePrice = input.purchasePrice;
    this.purchasedUPC = input.purchasedUPC;
    this.status = input.status = 'unknown';
    this.statusExplanation = input.statusExplanation = '';
    this.validMeasurementUnitID = input.validMeasurementUnitID = '';
    this.belongsToMealPlan = input.belongsToMealPlan = '';
  }
}
