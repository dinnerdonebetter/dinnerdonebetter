/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ValidIngredient } from './ValidIngredient';
import type { ValidMeasurementUnit } from './ValidMeasurementUnit';
export type ValidMeasurementUnitConversion = {
  archivedAt?: string;
  createdAt?: string;
  from?: ValidMeasurementUnit;
  id?: string;
  lastUpdatedAt?: string;
  modifier?: number;
  notes?: string;
  onlyForIngredient?: ValidIngredient;
  to?: ValidMeasurementUnit;
};
