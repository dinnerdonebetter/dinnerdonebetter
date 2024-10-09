// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IMealPlanOptionCreationRequestInput {
   mealScale: number;
 notes: string;
 assignedCook?: string;
 assignedDishwasher?: string;
 mealID: string;

}

export class MealPlanOptionCreationRequestInput implements IMealPlanOptionCreationRequestInput {
   mealScale: number;
 notes: string;
 assignedCook?: string;
 assignedDishwasher?: string;
 mealID: string;
constructor(input: Partial<MealPlanOptionCreationRequestInput> = {}) {
	 this.mealScale = input.mealScale = 0;
 this.notes = input.notes = '';
 this.assignedCook = input.assignedCook;
 this.assignedDishwasher = input.assignedDishwasher;
 this.mealID = input.mealID = '';
}
}