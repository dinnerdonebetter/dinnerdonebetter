/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { MealComponentCreationRequestInput } from './MealComponentCreationRequestInput';
export type MealCreationRequestInput = {
  components?: Array<MealComponentCreationRequestInput>;
  description?: string;
  eligibleForMealPlans?: boolean;
  maximumEstimatedPortions?: number;
  minimumEstimatedPortions?: number;
  name?: string;
};
