// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanEventUpdateRequestInput {
  endsAt?: string;
  mealName?: string;
  notes?: string;
  startsAt?: string;
}

export class MealPlanEventUpdateRequestInput implements IMealPlanEventUpdateRequestInput {
  endsAt?: string;
  mealName?: string;
  notes?: string;
  startsAt?: string;
  constructor(input: Partial<MealPlanEventUpdateRequestInput> = {}) {
    this.endsAt = input.endsAt;
    this.mealName = input.mealName;
    this.notes = input.notes;
    this.startsAt = input.startsAt;
  }
}
