// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMealPlanGroceryListItemStatus } from './enums';
import { OptionalNumberRange } from './number_range';

export interface IMealPlanGroceryListItemUpdateRequestInput {
  purchasedMeasurementUnitID?: string;
  quantityNeeded: OptionalNumberRange;
  status?: ValidMealPlanGroceryListItemStatus;
  validIngredientID?: string;
  validMeasurementUnitID?: string;
  belongsToMealPlan?: string;
  purchasePrice?: number;
  purchasedUPC?: string;
  quantityPurchased?: number;
  statusExplanation?: string;
}

export class MealPlanGroceryListItemUpdateRequestInput implements IMealPlanGroceryListItemUpdateRequestInput {
  purchasedMeasurementUnitID?: string;
  quantityNeeded: OptionalNumberRange;
  status?: ValidMealPlanGroceryListItemStatus;
  validIngredientID?: string;
  validMeasurementUnitID?: string;
  belongsToMealPlan?: string;
  purchasePrice?: number;
  purchasedUPC?: string;
  quantityPurchased?: number;
  statusExplanation?: string;
  constructor(input: Partial<MealPlanGroceryListItemUpdateRequestInput> = {}) {
    this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
    this.quantityNeeded = input.quantityNeeded = {};
    this.status = input.status;
    this.validIngredientID = input.validIngredientID;
    this.validMeasurementUnitID = input.validMeasurementUnitID;
    this.belongsToMealPlan = input.belongsToMealPlan;
    this.purchasePrice = input.purchasePrice;
    this.purchasedUPC = input.purchasedUPC;
    this.quantityPurchased = input.quantityPurchased;
    this.statusExplanation = input.statusExplanation;
  }
}
