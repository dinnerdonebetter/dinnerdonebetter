/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { MealPlanOption } from './MealPlanOption';
import type { RecipePrepTask } from './RecipePrepTask';
export type MealPlanTask = {
  assignedToUser?: string;
  completedAt?: string;
  createdAt?: string;
  creationExplanation?: string;
  id?: string;
  lastUpdatedAt?: string;
  mealPlanOption?: MealPlanOption;
  recipePrepTask?: RecipePrepTask;
  status?: string;
  statusExplanation?: string;
};
