// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanUpdateRequestInput {
  notes: string;
  votingDeadline: string;
}

export class MealPlanUpdateRequestInput implements IMealPlanUpdateRequestInput {
  notes: string;
  votingDeadline: string;
  constructor(input: Partial<MealPlanUpdateRequestInput> = {}) {
    this.notes = input.notes || '';
    this.votingDeadline = input.votingDeadline || '';
  }
}
