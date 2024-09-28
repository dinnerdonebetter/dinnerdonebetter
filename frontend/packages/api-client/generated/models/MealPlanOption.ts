/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Meal } from './Meal';
import type { MealPlanOptionVote } from './MealPlanOptionVote';
export type MealPlanOption = {
  archivedAt?: string;
  assignedCook?: string;
  assignedDishwasher?: string;
  belongsToMealPlanEvent?: string;
  chosen?: boolean;
  createdAt?: string;
  id?: string;
  lastUpdatedAt?: string;
  meal?: Meal;
  mealScale?: number;
  notes?: string;
  tieBroken?: boolean;
  votes?: Array<MealPlanOptionVote>;
};
