// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanEventCreationRequestInput } from './MealPlanEventCreationRequestInput';
import { ValidMealPlanElectionMethod } from './enums';

export interface IMealPlanCreationRequestInput {
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEventCreationRequestInput;
  notes: string;
  votingDeadline: string;
}

export class MealPlanCreationRequestInput implements IMealPlanCreationRequestInput {
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEventCreationRequestInput;
  notes: string;
  votingDeadline: string;
  constructor(input: Partial<MealPlanCreationRequestInput> = {}) {
    this.electionMethod = input.electionMethod = 'schulze';
    this.events = input.events = new MealPlanEventCreationRequestInput();
    this.notes = input.notes = '';
    this.votingDeadline = input.votingDeadline = '';
  }
}
