// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionVoteCreationInput {
  belongsToMealPlanOption: string;
  notes: string;
  rank: number;
  abstain: boolean;
}

export class MealPlanOptionVoteCreationInput implements IMealPlanOptionVoteCreationInput {
  belongsToMealPlanOption: string;
  notes: string;
  rank: number;
  abstain: boolean;
  constructor(input: Partial<MealPlanOptionVoteCreationInput> = {}) {
    this.belongsToMealPlanOption = input.belongsToMealPlanOption = '';
    this.notes = input.notes = '';
    this.rank = input.rank = 0;
    this.abstain = input.abstain = false;
  }
}
