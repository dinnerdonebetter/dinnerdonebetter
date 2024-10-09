// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanTaskStatus } from './enums';

export interface IMealPlanTaskStatusChangeRequestInput {
  assignedToUser: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
}

export class MealPlanTaskStatusChangeRequestInput implements IMealPlanTaskStatusChangeRequestInput {
  assignedToUser: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  constructor(input: Partial<MealPlanTaskStatusChangeRequestInput> = {}) {
    this.assignedToUser = input.assignedToUser || '';
    this.status = input.status || 'unfinished';
    this.statusExplanation = input.statusExplanation || '';
  }
}
