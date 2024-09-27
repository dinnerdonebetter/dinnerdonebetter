/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ValidIngredient } from './ValidIngredient';
import type { ValidMeasurementUnit } from './ValidMeasurementUnit';
export type ValidIngredientMeasurementUnit = {
  archivedAt?: string;
  createdAt?: string;
  id?: string;
  ingredient?: ValidIngredient;
  lastUpdatedAt?: string;
  maximumAllowableQuantity?: number;
  measurementUnit?: ValidMeasurementUnit;
  minimumAllowableQuantity?: number;
  notes?: string;
};
