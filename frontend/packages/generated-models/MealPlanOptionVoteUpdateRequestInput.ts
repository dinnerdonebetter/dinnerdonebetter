// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionVoteUpdateRequestInput {
  belongsToMealPlanOption: string;
  notes?: string;
  rank?: number;
  abstain?: boolean;
}

export class MealPlanOptionVoteUpdateRequestInput implements IMealPlanOptionVoteUpdateRequestInput {
  belongsToMealPlanOption: string;
  notes?: string;
  rank?: number;
  abstain?: boolean;
  constructor(input: Partial<MealPlanOptionVoteUpdateRequestInput> = {}) {
    this.belongsToMealPlanOption = input.belongsToMealPlanOption = '';
    this.notes = input.notes;
    this.rank = input.rank;
    this.abstain = input.abstain;
  }
}
