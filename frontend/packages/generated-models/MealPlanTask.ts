// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';
import { RecipePrepTask } from './RecipePrepTask';
import { MealPlanTaskStatus } from './enums';

export interface IMealPlanTask {
  assignedToUser?: string;
  createdAt: string;
  creationExplanation: string;
  id: string;
  lastUpdatedAt?: string;
  statusExplanation: string;
  completedAt?: string;
  mealPlanOption: MealPlanOption;
  recipePrepTask: RecipePrepTask;
  status: MealPlanTaskStatus;
}

export class MealPlanTask implements IMealPlanTask {
  assignedToUser?: string;
  createdAt: string;
  creationExplanation: string;
  id: string;
  lastUpdatedAt?: string;
  statusExplanation: string;
  completedAt?: string;
  mealPlanOption: MealPlanOption;
  recipePrepTask: RecipePrepTask;
  status: MealPlanTaskStatus;
  constructor(input: Partial<MealPlanTask> = {}) {
    this.assignedToUser = input.assignedToUser;
    this.createdAt = input.createdAt = '';
    this.creationExplanation = input.creationExplanation = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.statusExplanation = input.statusExplanation = '';
    this.completedAt = input.completedAt;
    this.mealPlanOption = input.mealPlanOption = new MealPlanOption();
    this.recipePrepTask = input.recipePrepTask = new RecipePrepTask();
    this.status = input.status = 'unfinished';
  }
}
