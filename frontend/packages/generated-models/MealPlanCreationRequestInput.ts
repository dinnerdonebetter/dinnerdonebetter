// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanEventCreationRequestInput } from './MealPlanEventCreationRequestInput';
import { ValidMealPlanElectionMethod } from './enums';

export interface IMealPlanCreationRequestInput {
  events: MealPlanEventCreationRequestInput;
  notes: string;
  votingDeadline: string;
  electionMethod: ValidMealPlanElectionMethod;
}

export class MealPlanCreationRequestInput implements IMealPlanCreationRequestInput {
  events: MealPlanEventCreationRequestInput;
  notes: string;
  votingDeadline: string;
  electionMethod: ValidMealPlanElectionMethod;
  constructor(input: Partial<MealPlanCreationRequestInput> = {}) {
    this.events = input.events = new MealPlanEventCreationRequestInput();
    this.notes = input.notes = '';
    this.votingDeadline = input.votingDeadline = '';
    this.electionMethod = input.electionMethod = 'schulze';
  }
}
