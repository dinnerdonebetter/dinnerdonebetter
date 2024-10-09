// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { OptionalNumberRange } from './number_range';

export interface IMealPlanGroceryListItemUpdateRequestInput {
  validMeasurementUnitID?: string;
  belongsToMealPlan?: string;
  purchasePrice?: number;
  purchasedMeasurementUnitID?: string;
  purchasedUPC?: string;
  validIngredientID?: string;
  quantityNeeded: OptionalNumberRange;
  quantityPurchased?: number;
  status?: ValidMealPlanGroceryListItemStatus;
  statusExplanation?: string;
}

export class MealPlanGroceryListItemUpdateRequestInput implements IMealPlanGroceryListItemUpdateRequestInput {
  validMeasurementUnitID?: string;
  belongsToMealPlan?: string;
  purchasePrice?: number;
  purchasedMeasurementUnitID?: string;
  purchasedUPC?: string;
  validIngredientID?: string;
  quantityNeeded: OptionalNumberRange;
  quantityPurchased?: number;
  status?: ValidMealPlanGroceryListItemStatus;
  statusExplanation?: string;
  constructor(input: Partial<MealPlanGroceryListItemUpdateRequestInput> = {}) {
    this.validMeasurementUnitID = input.validMeasurementUnitID;
    this.belongsToMealPlan = input.belongsToMealPlan;
    this.purchasePrice = input.purchasePrice;
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
    this.purchasedUPC = input.purchasedUPC;
    this.validIngredientID = input.validIngredientID;
    this.quantityNeeded = input.quantityNeeded = {};
    this.quantityPurchased = input.quantityPurchased;
    this.status = input.status;
    this.statusExplanation = input.statusExplanation;
  }
}
