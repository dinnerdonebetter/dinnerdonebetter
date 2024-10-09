// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidMealPlanGroceryListItemStatus } from './enums';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IMealPlanGroceryListItemCreationRequestInput {
   belongsToMealPlan: string;
 purchasePrice?: number;
 quantityPurchased?: number;
 status: ValidMealPlanGroceryListItemStatus;
 validMeasurementUnitID: string;
 purchasedMeasurementUnitID?: string;
 purchasedUPC?: string;
 quantityNeeded: NumberRangeWithOptionalMax;
 statusExplanation: string;
 validIngredientID: string;

}

export class MealPlanGroceryListItemCreationRequestInput implements IMealPlanGroceryListItemCreationRequestInput {
   belongsToMealPlan: string;
 purchasePrice?: number;
 quantityPurchased?: number;
 status: ValidMealPlanGroceryListItemStatus;
 validMeasurementUnitID: string;
 purchasedMeasurementUnitID?: string;
 purchasedUPC?: string;
 quantityNeeded: NumberRangeWithOptionalMax;
 statusExplanation: string;
 validIngredientID: string;
constructor(input: Partial<MealPlanGroceryListItemCreationRequestInput> = {}) {
	 this.belongsToMealPlan = input.belongsToMealPlan = '';
 this.purchasePrice = input.purchasePrice;
 this.quantityPurchased = input.quantityPurchased;
 this.status = input.status = 'unknown';
 this.validMeasurementUnitID = input.validMeasurementUnitID = '';
 this.purchasedMeasurementUnitID = input.purchasedMeasurementUnitID;
 this.purchasedUPC = input.purchasedUPC;
 this.quantityNeeded = input.quantityNeeded = { min: 0 };
 this.statusExplanation = input.statusExplanation = '';
 this.validIngredientID = input.validIngredientID = '';
}
}