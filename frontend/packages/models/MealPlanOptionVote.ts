// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IMealPlanOptionVote {
  createdAt: string;
  rank: number;
  archivedAt: string;
  belongsToMealPlanOption: string;
  byUser: string;
  notes: string;
  abstain: boolean;
  id: string;
  lastUpdatedAt: string;
}

export class MealPlanOptionVote implements IMealPlanOptionVote {
  createdAt: string;
  rank: number;
  archivedAt: string;
  belongsToMealPlanOption: string;
  byUser: string;
  notes: string;
  abstain: boolean;
  id: string;
  lastUpdatedAt: string;
  constructor(input: Partial<MealPlanOptionVote> = {}) {
    this.createdAt = input.createdAt || '';
    this.rank = input.rank || 0;
    this.archivedAt = input.archivedAt || '';
    this.belongsToMealPlanOption = input.belongsToMealPlanOption || '';
    this.byUser = input.byUser || '';
    this.notes = input.notes || '';
    this.abstain = input.abstain || false;
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
  }
}
