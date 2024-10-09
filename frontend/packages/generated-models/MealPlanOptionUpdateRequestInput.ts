// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionUpdateRequestInput {
  mealID?: string;
  mealScale?: number;
  notes?: string;
  assignedCook?: string;
  assignedDishwasher?: string;
}

export class MealPlanOptionUpdateRequestInput implements IMealPlanOptionUpdateRequestInput {
  mealID?: string;
  mealScale?: number;
  notes?: string;
  assignedCook?: string;
  assignedDishwasher?: string;
  constructor(input: Partial<MealPlanOptionUpdateRequestInput> = {}) {
    this.mealID = input.mealID;
    this.mealScale = input.mealScale;
    this.notes = input.notes;
    this.assignedCook = input.assignedCook;
    this.assignedDishwasher = input.assignedDishwasher;
  }
}
