/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ValidInstrument } from './ValidInstrument';
export type RecipeStepInstrument = {
  archivedAt?: string;
  belongsToRecipeStep?: string;
  createdAt?: string;
  id?: string;
  instrument?: ValidInstrument;
  lastUpdatedAt?: string;
  maximumQuantity?: number;
  minimumQuantity?: number;
  name?: string;
  notes?: string;
  optionIndex?: number;
  optional?: boolean;
  preferenceRank?: number;
  recipeStepProductID?: string;
};
