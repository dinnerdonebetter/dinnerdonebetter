// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOptionCreationRequestInput } from './MealPlanOptionCreationRequestInput';

export interface IMealPlanEventCreationRequestInput {
  endsAt: string;
  mealName: string;
  notes: string;
  options: MealPlanOptionCreationRequestInput[];
  startsAt: string;
}

export class MealPlanEventCreationRequestInput implements IMealPlanEventCreationRequestInput {
  endsAt: string;
  mealName: string;
  notes: string;
  options: MealPlanOptionCreationRequestInput[];
  startsAt: string;
  constructor(input: Partial<MealPlanEventCreationRequestInput> = {}) {
    this.endsAt = input.endsAt || '';
    this.mealName = input.mealName || '';
    this.notes = input.notes || '';
    this.options = input.options || [];
    this.startsAt = input.startsAt || '';
  }
}
