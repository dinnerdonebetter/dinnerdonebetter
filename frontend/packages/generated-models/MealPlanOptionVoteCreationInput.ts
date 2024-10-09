// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionVoteCreationInput {
  notes: string;
  rank: number;
  abstain: boolean;
  belongsToMealPlanOption: string;
}

export class MealPlanOptionVoteCreationInput implements IMealPlanOptionVoteCreationInput {
  notes: string;
  rank: number;
  abstain: boolean;
  belongsToMealPlanOption: string;
  constructor(input: Partial<MealPlanOptionVoteCreationInput> = {}) {
    this.notes = input.notes = '';
    this.rank = input.rank = 0;
    this.abstain = input.abstain = false;
    this.belongsToMealPlanOption = input.belongsToMealPlanOption = '';
  }
}
