// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanTaskStatus } from './enums';

export interface IMealPlanTaskCreationRequestInput {
  assignedToUser?: string;
  creationExplanation: string;
  mealPlanOptionID: string;
  recipePrepTaskID: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
}

export class MealPlanTaskCreationRequestInput implements IMealPlanTaskCreationRequestInput {
  assignedToUser?: string;
  creationExplanation: string;
  mealPlanOptionID: string;
  recipePrepTaskID: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  constructor(input: Partial<MealPlanTaskCreationRequestInput> = {}) {
    this.assignedToUser = input.assignedToUser;
    this.creationExplanation = input.creationExplanation = '';
    this.mealPlanOptionID = input.mealPlanOptionID = '';
    this.recipePrepTaskID = input.recipePrepTaskID = '';
    this.status = input.status = 'unfinished';
    this.statusExplanation = input.statusExplanation = '';
  }
}
