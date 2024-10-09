// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidMealPlanGroceryListItemStatus } from './enums';
 import { OptionalNumberRange } from './number_range';


export interface IMealPlanGroceryListItemUpdateRequestInput {
   statusExplanation?: string;
 belongsToMealPlan?: string;
 quantityNeeded: OptionalNumberRange;
 purchasedUPC?: string;
 quantityPurchased?: number;
 status?: ValidMealPlanGroceryListItemStatus;
 validIngredientID?: string;
 validMeasurementUnitID?: string;
 purchasePrice?: number;
 purchasedMeasurementUnitID?: string;

}

export class MealPlanGroceryListItemUpdateRequestInput implements IMealPlanGroceryListItemUpdateRequestInput {
   statusExplanation?: string;
 belongsToMealPlan?: string;
 quantityNeeded: OptionalNumberRange;
 purchasedUPC?: string;
 quantityPurchased?: number;
 status?: ValidMealPlanGroceryListItemStatus;
 validIngredientID?: string;
 validMeasurementUnitID?: string;
 purchasePrice?: number;
 purchasedMeasurementUnitID?: string;
constructor(input: Partial<MealPlanGroceryListItemUpdateRequestInput> = {}) {
	 this.statusExplanation = input.statusExplanation;
 this.belongsToMealPlan = input.belongsToMealPlan;
 this.quantityNeeded = input.quantityNeeded = {};
 this.purchasedUPC = input.purchasedUPC;
 this.quantityPurchased = input.quantityPurchased;
 this.status = input.status;
 this.validIngredientID = input.validIngredientID;
 this.validMeasurementUnitID = input.validMeasurementUnitID;
 this.purchasePrice = input.purchasePrice;
 this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
}
}