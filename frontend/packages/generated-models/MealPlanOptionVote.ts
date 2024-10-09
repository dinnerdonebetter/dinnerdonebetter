// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionVote {
  id: string;
  lastUpdatedAt?: string;
  notes: string;
  rank: number;
  abstain: boolean;
  archivedAt?: string;
  belongsToMealPlanOption: string;
  byUser: string;
  createdAt: string;
}

export class MealPlanOptionVote implements IMealPlanOptionVote {
  id: string;
  lastUpdatedAt?: string;
  notes: string;
  rank: number;
  abstain: boolean;
  archivedAt?: string;
  belongsToMealPlanOption: string;
  byUser: string;
  createdAt: string;
  constructor(input: Partial<MealPlanOptionVote> = {}) {
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.rank = input.rank = 0;
    this.abstain = input.abstain = false;
    this.archivedAt = input.archivedAt;
    this.belongsToMealPlanOption = input.belongsToMealPlanOption = '';
    this.byUser = input.byUser = '';
    this.createdAt = input.createdAt = '';
  }
}
