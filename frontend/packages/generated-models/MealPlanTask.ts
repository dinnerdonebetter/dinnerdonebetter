// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';
import { RecipePrepTask } from './RecipePrepTask';
import { MealPlanTaskStatus } from './enums';

export interface IMealPlanTask {
  createdAt: string;
  lastUpdatedAt?: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  assignedToUser?: string;
  completedAt?: string;
  mealPlanOption: MealPlanOption;
  recipePrepTask: RecipePrepTask;
  creationExplanation: string;
  id: string;
}

export class MealPlanTask implements IMealPlanTask {
  createdAt: string;
  lastUpdatedAt?: string;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  assignedToUser?: string;
  completedAt?: string;
  mealPlanOption: MealPlanOption;
  recipePrepTask: RecipePrepTask;
  creationExplanation: string;
  id: string;
  constructor(input: Partial<MealPlanTask> = {}) {
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.status = input.status = 'unfinished';
    this.statusExplanation = input.statusExplanation = '';
    this.assignedToUser = input.assignedToUser;
    this.completedAt = input.completedAt;
    this.mealPlanOption = input.mealPlanOption = new MealPlanOption();
    this.recipePrepTask = input.recipePrepTask = new RecipePrepTask();
    this.creationExplanation = input.creationExplanation = '';
    this.id = input.id = '';
  }
}
