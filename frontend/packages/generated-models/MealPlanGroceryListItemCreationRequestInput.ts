// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealPlanGroceryListItemCreationRequestInput {
  validMeasurementUnitID: string;
  belongsToMealPlan: string;
  purchasePrice?: number;
  quantityPurchased?: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  validIngredientID: string;
  purchasedMeasurementUnitID?: string;
  purchasedUPC?: string;
  quantityNeeded: NumberRangeWithOptionalMax;
}

export class MealPlanGroceryListItemCreationRequestInput implements IMealPlanGroceryListItemCreationRequestInput {
  validMeasurementUnitID: string;
  belongsToMealPlan: string;
  purchasePrice?: number;
  quantityPurchased?: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  validIngredientID: string;
  purchasedMeasurementUnitID?: string;
  purchasedUPC?: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  constructor(input: Partial<MealPlanGroceryListItemCreationRequestInput> = {}) {
    this.validMeasurementUnitID = input.validMeasurementUnitID = '';
    this.belongsToMealPlan = input.belongsToMealPlan = '';
    this.purchasePrice = input.purchasePrice;
    this.quantityPurchased = input.quantityPurchased;
    this.status = input.status = 'unknown';
    this.statusExplanation = input.statusExplanation = '';
    this.validIngredientID = input.validIngredientID = '';
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
    this.purchasedUPC = input.purchasedUPC;
    this.quantityNeeded = input.quantityNeeded = { min: 0 };
  }
}
