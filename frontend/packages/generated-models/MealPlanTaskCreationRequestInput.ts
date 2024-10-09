// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanTaskStatus } from './enums';

export interface IMealPlanTaskCreationRequestInput {
  creationExplanation: string;
  mealPlanOptionID: string;
  recipePrepTaskID: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  assignedToUser?: string;
}

export class MealPlanTaskCreationRequestInput implements IMealPlanTaskCreationRequestInput {
  creationExplanation: string;
  mealPlanOptionID: string;
  recipePrepTaskID: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  assignedToUser?: string;
  constructor(input: Partial<MealPlanTaskCreationRequestInput> = {}) {
    this.creationExplanation = input.creationExplanation = '';
    this.mealPlanOptionID = input.mealPlanOptionID = '';
    this.recipePrepTaskID = input.recipePrepTaskID = '';
    this.status = input.status = 'unfinished';
    this.statusExplanation = input.statusExplanation = '';
    this.assignedToUser = input.assignedToUser;
  }
}
