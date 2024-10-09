// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IMealPlanOptionVote {
   lastUpdatedAt?: string;
 notes: string;
 rank: number;
 abstain: boolean;
 belongsToMealPlanOption: string;
 byUser: string;
 archivedAt?: string;
 createdAt: string;
 id: string;

}

export class MealPlanOptionVote implements IMealPlanOptionVote {
   lastUpdatedAt?: string;
 notes: string;
 rank: number;
 abstain: boolean;
 belongsToMealPlanOption: string;
 byUser: string;
 archivedAt?: string;
 createdAt: string;
 id: string;
constructor(input: Partial<MealPlanOptionVote> = {}) {
	 this.lastUpdatedAt = input.lastUpdatedAt;
 this.notes = input.notes = '';
 this.rank = input.rank = 0;
 this.abstain = input.abstain = false;
 this.belongsToMealPlanOption = input.belongsToMealPlanOption = '';
 this.byUser = input.byUser = '';
 this.archivedAt = input.archivedAt;
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
}
}