/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ValidMeasurementUnit } from './ValidMeasurementUnit';
export type RecipeStepProduct = {
  archivedAt?: string;
  belongsToRecipeStep?: string;
  compostable?: boolean;
  containedInVesselIndex?: number;
  createdAt?: string;
  id?: string;
  index?: number;
  isLiquid?: boolean;
  isWaste?: boolean;
  lastUpdatedAt?: string;
  maximumQuantity?: number;
  maximumStorageDurationInSeconds?: number;
  maximumStorageTemperatureInCelsius?: number;
  measurementUnit?: ValidMeasurementUnit;
  minimumQuantity?: number;
  minimumStorageTemperatureInCelsius?: number;
  name?: string;
  quantityNotes?: string;
  storageInstructions?: string;
  type?: string;
};
