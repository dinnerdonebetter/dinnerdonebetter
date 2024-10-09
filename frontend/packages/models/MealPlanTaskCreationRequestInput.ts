// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanTaskStatus } from './enums';

export interface IMealPlanTaskCreationRequestInput {
  mealPlanOptionID: string;
  recipePrepTaskID: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  assignedToUser: string;
  creationExplanation: string;
}

export class MealPlanTaskCreationRequestInput implements IMealPlanTaskCreationRequestInput {
  mealPlanOptionID: string;
  recipePrepTaskID: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  assignedToUser: string;
  creationExplanation: string;
  constructor(input: Partial<MealPlanTaskCreationRequestInput> = {}) {
    this.mealPlanOptionID = input.mealPlanOptionID || '';
    this.recipePrepTaskID = input.recipePrepTaskID || '';
    this.status = input.status || 'unfinished';
    this.statusExplanation = input.statusExplanation || '';
    this.assignedToUser = input.assignedToUser || '';
    this.creationExplanation = input.creationExplanation || '';
  }
}
