// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanEventCreationRequestInput } from './MealPlanEventCreationRequestInput.gen';
import { ValidMealPlanElectionMethod } from './enums.gen';

export interface IMealPlanCreationRequestInput {
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEventCreationRequestInput[];
  notes: string;
  votingDeadline: string;
}

export class MealPlanCreationRequestInput implements IMealPlanCreationRequestInput {
  electionMethod: ValidMealPlanElectionMethod;
  events: MealPlanEventCreationRequestInput[];
  notes: string;
  votingDeadline: string;
  constructor(input: Partial<MealPlanCreationRequestInput> = {}) {
    this.electionMethod = input.electionMethod || 'schulze';
    this.events = input.events || [];
    this.notes = input.notes || '';
    this.votingDeadline = input.votingDeadline || '';
  }
}
