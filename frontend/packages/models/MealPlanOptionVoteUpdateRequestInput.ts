// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionVoteUpdateRequestInput {
  abstain: boolean;
  belongsToMealPlanOption: string;
  notes: string;
  rank: number;
}

export class MealPlanOptionVoteUpdateRequestInput implements IMealPlanOptionVoteUpdateRequestInput {
  abstain: boolean;
  belongsToMealPlanOption: string;
  notes: string;
  rank: number;
  constructor(input: Partial<MealPlanOptionVoteUpdateRequestInput> = {}) {
    this.abstain = input.abstain || false;
    this.belongsToMealPlanOption = input.belongsToMealPlanOption || '';
    this.notes = input.notes || '';
    this.rank = input.rank || 0;
  }
}
