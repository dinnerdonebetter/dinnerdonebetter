// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionUpdateRequestInput {
  mealScale?: number;
  notes?: string;
  assignedCook?: string;
  assignedDishwasher?: string;
  mealID?: string;
}

export class MealPlanOptionUpdateRequestInput implements IMealPlanOptionUpdateRequestInput {
  mealScale?: number;
  notes?: string;
  assignedCook?: string;
  assignedDishwasher?: string;
  mealID?: string;
  constructor(input: Partial<MealPlanOptionUpdateRequestInput> = {}) {
    this.mealScale = input.mealScale;
    this.notes = input.notes;
    this.assignedCook = input.assignedCook;
    this.assignedDishwasher = input.assignedDishwasher;
    this.mealID = input.mealID;
  }
}
