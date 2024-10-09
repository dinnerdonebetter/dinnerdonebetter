// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { OptionalNumberRange } from './number_range';

export interface IMealPlanGroceryListItemUpdateRequestInput {
  belongsToMealPlan: string;
  purchasedUPC: string;
  statusExplanation: string;
  validMeasurementUnitID: string;
  purchasePrice: number;
  purchasedMeasurementUnitID: string;
  quantityNeeded: OptionalNumberRange;
  quantityPurchased: number;
  status: ValidMealPlanGroceryListItemStatus;
  validIngredientID: string;
}

export class MealPlanGroceryListItemUpdateRequestInput implements IMealPlanGroceryListItemUpdateRequestInput {
  belongsToMealPlan: string;
  purchasedUPC: string;
  statusExplanation: string;
  validMeasurementUnitID: string;
  purchasePrice: number;
  purchasedMeasurementUnitID: string;
  quantityNeeded: OptionalNumberRange;
  quantityPurchased: number;
  status: ValidMealPlanGroceryListItemStatus;
  validIngredientID: string;
  constructor(input: Partial<MealPlanGroceryListItemUpdateRequestInput> = {}) {
    this.belongsToMealPlan = input.belongsToMealPlan || '';
    this.purchasedUPC = input.purchasedUPC || '';
    this.statusExplanation = input.statusExplanation || '';
    this.validMeasurementUnitID = input.validMeasurementUnitID || '';
    this.purchasePrice = input.purchasePrice || 0;
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID || '';
    this.quantityNeeded = input.quantityNeeded || {};
    this.quantityPurchased = input.quantityPurchased || 0;
    this.status = input.status || 'unknown';
    this.validIngredientID = input.validIngredientID || '';
  }
}
