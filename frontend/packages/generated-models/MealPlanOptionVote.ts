// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionVote {
  abstain: boolean;
  belongsToMealPlanOption: string;
  id: string;
  notes: string;
  rank: number;
  archivedAt?: string;
  byUser: string;
  createdAt: string;
  lastUpdatedAt?: string;
}

export class MealPlanOptionVote implements IMealPlanOptionVote {
  abstain: boolean;
  belongsToMealPlanOption: string;
  id: string;
  notes: string;
  rank: number;
  archivedAt?: string;
  byUser: string;
  createdAt: string;
  lastUpdatedAt?: string;
  constructor(input: Partial<MealPlanOptionVote> = {}) {
    this.abstain = input.abstain = false;
    this.belongsToMealPlanOption = input.belongsToMealPlanOption = '';
    this.id = input.id = '';
    this.notes = input.notes = '';
    this.rank = input.rank = 0;
    this.archivedAt = input.archivedAt;
    this.byUser = input.byUser = '';
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
