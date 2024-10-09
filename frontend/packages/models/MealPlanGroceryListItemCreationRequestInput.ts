// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { NumberRangeWithOptionalMax } from './number_range';

export interface IMealPlanGroceryListItemCreationRequestInput {
  status: ValidMealPlanGroceryListItemStatus;
  belongsToMealPlan: string;
  purchasedUPC: string;
  quantityPurchased: number;
  statusExplanation: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
  purchasePrice: number;
  purchasedMeasurementUnitID: string;
  quantityNeeded: NumberRangeWithOptionalMax;
}

export class MealPlanGroceryListItemCreationRequestInput implements IMealPlanGroceryListItemCreationRequestInput {
  status: ValidMealPlanGroceryListItemStatus;
  belongsToMealPlan: string;
  purchasedUPC: string;
  quantityPurchased: number;
  statusExplanation: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
  purchasePrice: number;
  purchasedMeasurementUnitID: string;
  quantityNeeded: NumberRangeWithOptionalMax;
  constructor(input: Partial<MealPlanGroceryListItemCreationRequestInput> = {}) {
    this.status = input.status || 'unknown';
    this.belongsToMealPlan = input.belongsToMealPlan || '';
    this.purchasedUPC = input.purchasedUPC || '';
    this.quantityPurchased = input.quantityPurchased || 0;
    this.statusExplanation = input.statusExplanation || '';
    this.validIngredientID = input.validIngredientID || '';
    this.validMeasurementUnitID = input.validMeasurementUnitID || '';
    this.purchasePrice = input.purchasePrice || 0;
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID || '';
    this.quantityNeeded = input.quantityNeeded || { min: 0 };
  }
}
