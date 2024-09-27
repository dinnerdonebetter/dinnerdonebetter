/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepUpdateRequestInput = {
  properties: {
    belongsToRecipe: {
      type: 'string',
    },
    conditionExpression: {
      type: 'string',
    },
    explicitInstructions: {
      type: 'string',
    },
    index: {
      type: 'number',
      format: 'int64',
    },
    maximumEstimatedTimeInSeconds: {
      type: 'number',
      format: 'int64',
    },
    maximumTemperatureInCelsius: {
      type: 'number',
      format: 'double',
    },
    minimumEstimatedTimeInSeconds: {
      type: 'number',
      format: 'int64',
    },
    minimumTemperatureInCelsius: {
      type: 'number',
      format: 'double',
    },
    notes: {
      type: 'string',
    },
    optional: {
      type: 'boolean',
    },
    preparation: {
      type: 'ValidPreparation',
    },
    startTimerAutomatically: {
      type: 'boolean',
    },
  },
} as const;
