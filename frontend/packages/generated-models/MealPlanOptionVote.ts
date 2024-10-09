// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionVote {
  abstain: boolean;
  archivedAt?: string;
  byUser: string;
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  rank: number;
  belongsToMealPlanOption: string;
  notes: string;
}

export class MealPlanOptionVote implements IMealPlanOptionVote {
  abstain: boolean;
  archivedAt?: string;
  byUser: string;
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  rank: number;
  belongsToMealPlanOption: string;
  notes: string;
  constructor(input: Partial<MealPlanOptionVote> = {}) {
    this.abstain = input.abstain = false;
    this.archivedAt = input.archivedAt;
    this.byUser = input.byUser = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.rank = input.rank = 0;
    this.belongsToMealPlanOption = input.belongsToMealPlanOption = '';
    this.notes = input.notes = '';
  }
}
