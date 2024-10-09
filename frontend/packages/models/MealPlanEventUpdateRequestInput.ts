// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanEventUpdateRequestInput {
  startsAt: string;
  endsAt: string;
  mealName: string;
  notes: string;
}

export class MealPlanEventUpdateRequestInput implements IMealPlanEventUpdateRequestInput {
  startsAt: string;
  endsAt: string;
  mealName: string;
  notes: string;
  constructor(input: Partial<MealPlanEventUpdateRequestInput> = {}) {
    this.startsAt = input.startsAt || '';
    this.endsAt = input.endsAt || '';
    this.mealName = input.mealName || '';
    this.notes = input.notes || '';
  }
}
