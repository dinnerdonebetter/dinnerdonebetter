/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { MealPlanEventCreationRequestInput } from './MealPlanEventCreationRequestInput';
export type MealPlanCreationRequestInput = {
  electionMethod?: string;
  events?: Array<MealPlanEventCreationRequestInput>;
  notes?: string;
  votingDeadline?: string;
};
