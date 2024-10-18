// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionUpdateRequestInput {
  assignedCook: string;
  assignedDishwasher: string;
  mealID: string;
  mealScale: number;
  notes: string;
}

export class MealPlanOptionUpdateRequestInput implements IMealPlanOptionUpdateRequestInput {
  assignedCook: string;
  assignedDishwasher: string;
  mealID: string;
  mealScale: number;
  notes: string;
  constructor(input: Partial<MealPlanOptionUpdateRequestInput> = {}) {
    this.assignedCook = input.assignedCook || '';
    this.assignedDishwasher = input.assignedDishwasher || '';
    this.mealID = input.mealID || '';
    this.mealScale = input.mealScale || 0;
    this.notes = input.notes || '';
  }
}
