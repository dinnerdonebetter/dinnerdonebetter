/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepProductUpdateRequestInput = {
  properties: {
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
    measurementUnitID: {
      type: 'string',
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
