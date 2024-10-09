// GENERATED CODE, DO NOT EDIT MANUALLY

 import { MealPlanEventCreationRequestInput } from './MealPlanEventCreationRequestInput';
 import { ValidMealPlanElectionMethod } from './enums';


export interface IMealPlanCreationRequestInput {
   votingDeadline: string;
 electionMethod: ValidMealPlanElectionMethod;
 events: MealPlanEventCreationRequestInput;
 notes: string;

}

export class MealPlanCreationRequestInput implements IMealPlanCreationRequestInput {
   votingDeadline: string;
 electionMethod: ValidMealPlanElectionMethod;
 events: MealPlanEventCreationRequestInput;
 notes: string;
constructor(input: Partial<MealPlanCreationRequestInput> = {}) {
	 this.votingDeadline = input.votingDeadline = '';
 this.electionMethod = input.electionMethod = 'schulze';
 this.events = input.events = new MealPlanEventCreationRequestInput();
 this.notes = input.notes = '';
}
}