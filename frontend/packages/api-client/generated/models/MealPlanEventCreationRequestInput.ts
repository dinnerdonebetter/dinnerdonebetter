/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { MealPlanOptionCreationRequestInput } from './MealPlanOptionCreationRequestInput';
export type MealPlanEventCreationRequestInput = {
  endsAt?: string;
  mealName?: string;
  notes?: string;
  options?: Array<MealPlanOptionCreationRequestInput>;
  startsAt?: string;
};
