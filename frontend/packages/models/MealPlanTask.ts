// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';
import { RecipePrepTask } from './RecipePrepTask';
import { MealPlanTaskStatus } from './enums';

export interface IMealPlanTask {
  statusExplanation: string;
  id: string;
  recipePrepTask: RecipePrepTask;
  status: MealPlanTaskStatus;
  creationExplanation: string;
  lastUpdatedAt: string;
  mealPlanOption: MealPlanOption;
  assignedToUser: string;
  completedAt: string;
  createdAt: string;
}

export class MealPlanTask implements IMealPlanTask {
  statusExplanation: string;
  id: string;
  recipePrepTask: RecipePrepTask;
  status: MealPlanTaskStatus;
  creationExplanation: string;
  lastUpdatedAt: string;
  mealPlanOption: MealPlanOption;
  assignedToUser: string;
  completedAt: string;
  createdAt: string;
  constructor(input: Partial<MealPlanTask> = {}) {
    this.statusExplanation = input.statusExplanation || '';
    this.id = input.id || '';
    this.recipePrepTask = input.recipePrepTask || new RecipePrepTask();
    this.status = input.status || 'unfinished';
    this.creationExplanation = input.creationExplanation || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.mealPlanOption = input.mealPlanOption || new MealPlanOption();
    this.assignedToUser = input.assignedToUser || '';
    this.completedAt = input.completedAt || '';
    this.createdAt = input.createdAt || '';
  }
}
