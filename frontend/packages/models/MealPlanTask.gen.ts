// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption.gen';
import { RecipePrepTask } from './RecipePrepTask.gen';
import { MealPlanTaskStatus } from './enums.gen';

export interface IMealPlanTask {
  assignedToUser: string;
  completedAt: string;
  createdAt: string;
  creationExplanation: string;
  id: string;
  lastUpdatedAt: string;
  mealPlanOption: MealPlanOption;
  recipePrepTask: RecipePrepTask;
  status: MealPlanTaskStatus;
  statusExplanation: string;
}

export class MealPlanTask implements IMealPlanTask {
  assignedToUser: string;
  completedAt: string;
  createdAt: string;
  creationExplanation: string;
  id: string;
  lastUpdatedAt: string;
  mealPlanOption: MealPlanOption;
  recipePrepTask: RecipePrepTask;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  constructor(input: Partial<MealPlanTask> = {}) {
    this.assignedToUser = input.assignedToUser || '';
    this.completedAt = input.completedAt || '';
    this.createdAt = input.createdAt || '';
    this.creationExplanation = input.creationExplanation || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.mealPlanOption = input.mealPlanOption || new MealPlanOption();
    this.recipePrepTask = input.recipePrepTask || new RecipePrepTask();
    this.status = input.status || 'unfinished';
    this.statusExplanation = input.statusExplanation || '';
  }
}
