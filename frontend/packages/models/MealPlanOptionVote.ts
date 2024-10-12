// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IMealPlanOptionVote {
   abstain: boolean;
 archivedAt: string;
 belongsToMealPlanOption: string;
 byUser: string;
 createdAt: string;
 id: string;
 lastUpdatedAt: string;
 notes: string;
 rank: number;

}

export class MealPlanOptionVote implements IMealPlanOptionVote {
   abstain: boolean;
 archivedAt: string;
 belongsToMealPlanOption: string;
 byUser: string;
 createdAt: string;
 id: string;
 lastUpdatedAt: string;
 notes: string;
 rank: number;
constructor(input: Partial<MealPlanOptionVote> = {}) {
	 this.abstain = input.abstain || false;
 this.archivedAt = input.archivedAt || '';
 this.belongsToMealPlanOption = input.belongsToMealPlanOption || '';
 this.byUser = input.byUser || '';
 this.createdAt = input.createdAt || '';
 this.id = input.id || '';
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.notes = input.notes || '';
 this.rank = input.rank || 0;
}
}