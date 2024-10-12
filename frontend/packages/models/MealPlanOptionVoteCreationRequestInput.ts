// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealPlanOptionVoteCreationInput } from './MealPlanOptionVoteCreationInput';


export interface IMealPlanOptionVoteCreationRequestInput {
   votes: MealPlanOptionVoteCreationInput[];

}

export class MealPlanOptionVoteCreationRequestInput implements IMealPlanOptionVoteCreationRequestInput {
   votes: MealPlanOptionVoteCreationInput[];
constructor(input: Partial<MealPlanOptionVoteCreationRequestInput> = {}) {
	 this.votes = input.votes || [];
}
}