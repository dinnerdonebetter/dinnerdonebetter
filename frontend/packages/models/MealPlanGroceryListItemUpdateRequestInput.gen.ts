// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { OptionalNumberRange } from './number_range';

export interface IMealPlanGroceryListItemUpdateRequestInput {
  belongsToMealPlan: string;
  purchasePrice: number;
  purchasedMeasurementUnitID: string;
  purchasedUPC: string;
  quantityNeeded: OptionalNumberRange;
  quantityPurchased: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
}

export class MealPlanGroceryListItemUpdateRequestInput implements IMealPlanGroceryListItemUpdateRequestInput {
  belongsToMealPlan: string;
  purchasePrice: number;
  purchasedMeasurementUnitID: string;
  purchasedUPC: string;
  quantityNeeded: OptionalNumberRange;
  quantityPurchased: number;
  status: ValidMealPlanGroceryListItemStatus;
  statusExplanation: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
  constructor(input: Partial<MealPlanGroceryListItemUpdateRequestInput> = {}) {
    this.belongsToMealPlan = input.belongsToMealPlan || '';
    this.purchasePrice = input.purchasePrice || 0;
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID || '';
    this.purchasedUPC = input.purchasedUPC || '';
    this.quantityNeeded = input.quantityNeeded || {};
    this.quantityPurchased = input.quantityPurchased || 0;
    this.status = input.status || 'unknown';
    this.statusExplanation = input.statusExplanation || '';
    this.validIngredientID = input.validIngredientID || '';
    this.validMeasurementUnitID = input.validMeasurementUnitID || '';
  }
}
