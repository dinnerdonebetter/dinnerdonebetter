// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealPlanGroceryListItemCreationRequestInput {
  belongsToMealPlan: string;
  purchasePrice: number;
  purchasedMeasurementUnitID: string;
  purchasedUPC: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
}

export class MealPlanGroceryListItemCreationRequestInput implements IMealPlanGroceryListItemCreationRequestInput {
  belongsToMealPlan: string;
  purchasePrice: number;
  purchasedMeasurementUnitID: string;
  purchasedUPC: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  quantityPurchased: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
  constructor(input: Partial<MealPlanGroceryListItemCreationRequestInput> = {}) {
    this.belongsToMealPlan = input.belongsToMealPlan || '';
    this.purchasePrice = input.purchasePrice || 0;
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID || '';
    this.purchasedUPC = input.purchasedUPC || '';
    this.quantityNeeded = input.quantityNeeded || { min: 0 };
    this.quantityPurchased = input.quantityPurchased || 0;
    this.status = input.status || 'unknown';
    this.statusExplanation = input.statusExplanation || '';
    this.validIngredientID = input.validIngredientID || '';
    this.validMeasurementUnitID = input.validMeasurementUnitID || '';
  }
}
