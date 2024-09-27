/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepProduct = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToRecipeStep: {
      type: 'string',
    },
    compostable: {
      type: 'boolean',
    },
    containedInVesselIndex: {
      type: 'number',
      format: 'int64',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    id: {
      type: 'string',
    },
    index: {
      type: 'number',
      format: 'int64',
    },
    isLiquid: {
      type: 'boolean',
    },
    isWaste: {
      type: 'boolean',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    maximumQuantity: {
      type: 'number',
      format: 'double',
    },
    maximumStorageDurationInSeconds: {
      type: 'number',
      format: 'int64',
    },
    maximumStorageTemperatureInCelsius: {
      type: 'number',
      format: 'double',
    },
    measurementUnit: {
      type: 'ValidMeasurementUnit',
    },
    minimumQuantity: {
      type: 'number',
      format: 'double',
    },
    minimumStorageTemperatureInCelsius: {
      type: 'number',
      format: 'double',
    },
    name: {
      type: 'string',
    },
    quantityNotes: {
      type: 'string',
    },
    storageInstructions: {
      type: 'string',
    },
    type: {
      type: 'string',
    },
  },
} as const;
