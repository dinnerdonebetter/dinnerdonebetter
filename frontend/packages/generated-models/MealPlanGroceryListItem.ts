// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidIngredient } from './ValidIngredient';
 import { ValidMeasurementUnit } from './ValidMeasurementUnit';
 import { ValidMealPlanGroceryListItemStatus } from './enums';
 import { NumberRangeWithOptionalMax } from './number_range';


export interface IMealPlanGroceryListItem {
   archivedAt?: string;
 lastUpdatedAt?: string;
 status: ValidMealPlanGroceryListItemStatus;
 statusExplanation: string;
 measurementUnit: ValidMeasurementUnit;
 quantityNeeded: NumberRangeWithOptionalMax;
 belongsToMealPlan: string;
 createdAt: string;
 ingredient: ValidIngredient;
 purchasedMeasurementUnit?: ValidMeasurementUnit;
 quantityPurchased?: number;
 id: string;
 purchasePrice?: number;
 purchasedUPC?: string;

}

export class MealPlanGroceryListItem implements IMealPlanGroceryListItem {
   archivedAt?: string;
 lastUpdatedAt?: string;
 status: ValidMealPlanGroceryListItemStatus;
 statusExplanation: string;
 measurementUnit: ValidMeasurementUnit;
 quantityNeeded: NumberRangeWithOptionalMax;
 belongsToMealPlan: string;
 createdAt: string;
 ingredient: ValidIngredient;
 purchasedMeasurementUnit?: ValidMeasurementUnit;
 quantityPurchased?: number;
 id: string;
 purchasePrice?: number;
 purchasedUPC?: string;
constructor(input: Partial<MealPlanGroceryListItem> = {}) {
	 this.archivedAt = input.archivedAt;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.status = input.status = 'unknown';
 this.statusExplanation = input.statusExplanation = '';
 this.measurementUnit = input.measurementUnit = new ValidMeasurementUnit();
 this.quantityNeeded = input.quantityNeeded = { min: 0 };
 this.belongsToMealPlan = input.belongsToMealPlan = '';
 this.createdAt = input.createdAt = '';
 this.ingredient = input.ingredient = new ValidIngredient();
 this.purchasedMeasurementUnit = input.purchasedMeasurementUnit;
 this.quantityPurchased = input.quantityPurchased;
 this.id = input.id = '';
 this.purchasePrice = input.purchasePrice;
 this.purchasedUPC = input.purchasedUPC;
}
}