// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { OptionalNumberRange } from './number_range';

export interface IMealPlanGroceryListItemUpdateRequestInput {
  status?: ValidMealPlanGroceryListItemStatus;
  validIngredientID?: string;
  validMeasurementUnitID?: string;
  purchasePrice?: number;
  purchasedUPC?: string;
  quantityNeeded: OptionalNumberRange;
  quantityPurchased?: number;
  belongsToMealPlan?: string;
  purchasedMeasurementUnitID?: string;
  statusExplanation?: string;
}

export class MealPlanGroceryListItemUpdateRequestInput implements IMealPlanGroceryListItemUpdateRequestInput {
  status?: ValidMealPlanGroceryListItemStatus;
  validIngredientID?: string;
  validMeasurementUnitID?: string;
  purchasePrice?: number;
  purchasedUPC?: string;
  quantityNeeded: OptionalNumberRange;
  quantityPurchased?: number;
  belongsToMealPlan?: string;
  purchasedMeasurementUnitID?: string;
  statusExplanation?: string;
  constructor(input: Partial<MealPlanGroceryListItemUpdateRequestInput> = {}) {
    this.status = input.status;
    this.validIngredientID = input.validIngredientID;
    this.validMeasurementUnitID = input.validMeasurementUnitID;
    this.purchasePrice = input.purchasePrice;
    this.purchasedUPC = input.purchasedUPC;
    this.quantityNeeded = input.quantityNeeded = {};
    this.quantityPurchased = input.quantityPurchased;
    this.belongsToMealPlan = input.belongsToMealPlan;
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
    this.statusExplanation = input.statusExplanation;
  }
}
