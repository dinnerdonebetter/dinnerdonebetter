/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { MealPlanEvent } from './MealPlanEvent';
export type MealPlan = {
  archivedAt?: string;
  belongsToHousehold?: string;
  createdAt?: string;
  createdBy?: string;
  electionMethod?: string;
  events?: Array<MealPlanEvent>;
  groceryListInitialized?: boolean;
  id?: string;
  lastUpdatedAt?: string;
  notes?: string;
  status?: string;
  tasksCreated?: boolean;
  votingDeadline?: string;
};
