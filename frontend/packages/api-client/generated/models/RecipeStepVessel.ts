/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ValidVessel } from './ValidVessel';
export type RecipeStepVessel = {
  archivedAt?: string;
  belongsToRecipeStep?: string;
  createdAt?: string;
  id?: string;
  lastUpdatedAt?: string;
  maximumQuantity?: number;
  minimumQuantity?: number;
  name?: string;
  notes?: string;
  recipeStepProductID?: string;
  unavailableAfterStep?: boolean;
  vessel?: ValidVessel;
  vesselPreposition?: string;
};
