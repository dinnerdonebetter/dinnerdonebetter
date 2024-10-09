// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IMealPlanOptionVoteCreationInput {
   abstain: boolean;
 belongsToMealPlanOption: string;
 notes: string;
 rank: number;

}

export class MealPlanOptionVoteCreationInput implements IMealPlanOptionVoteCreationInput {
   abstain: boolean;
 belongsToMealPlanOption: string;
 notes: string;
 rank: number;
constructor(input: Partial<MealPlanOptionVoteCreationInput> = {}) {
	 this.abstain = input.abstain = false;
 this.belongsToMealPlanOption = input.belongsToMealPlanOption = '';
 this.notes = input.notes = '';
 this.rank = input.rank = 0;
}
}