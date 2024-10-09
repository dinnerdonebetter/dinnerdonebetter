// GENERATED CODE, DO NOT EDIT MANUALLY

import { MealPlanOption } from './MealPlanOption';
import { RecipePrepTask } from './RecipePrepTask';
import { MealPlanTaskStatus } from './enums';

export interface IMealPlanTask {
  mealPlanOption: MealPlanOption;
  recipePrepTask: RecipePrepTask;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  assignedToUser?: string;
  createdAt: string;
  id: string;
  completedAt?: string;
  creationExplanation: string;
  lastUpdatedAt?: string;
}

export class MealPlanTask implements IMealPlanTask {
  mealPlanOption: MealPlanOption;
  recipePrepTask: RecipePrepTask;
  status: MealPlanTaskStatus;
  statusExplanation: string;
  assignedToUser?: string;
  createdAt: string;
  id: string;
  completedAt?: string;
  creationExplanation: string;
  lastUpdatedAt?: string;
  constructor(input: Partial<MealPlanTask> = {}) {
    this.mealPlanOption = input.mealPlanOption = new MealPlanOption();
    this.recipePrepTask = input.recipePrepTask = new RecipePrepTask();
    this.status = input.status = 'unfinished';
    this.statusExplanation = input.statusExplanation = '';
    this.assignedToUser = input.assignedToUser;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.completedAt = input.completedAt;
    this.creationExplanation = input.creationExplanation = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
