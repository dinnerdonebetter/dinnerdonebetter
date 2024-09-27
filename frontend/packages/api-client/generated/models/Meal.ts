/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { MealComponent } from './MealComponent';
export type Meal = {
  archivedAt?: string;
  components?: Array<MealComponent>;
  createdAt?: string;
  createdByUser?: string;
  description?: string;
  eligibleForMealPlans?: boolean;
  id?: string;
  lastUpdatedAt?: string;
  maximumEstimatedPortions?: number;
  minimumEstimatedPortions?: number;
  name?: string;
};
