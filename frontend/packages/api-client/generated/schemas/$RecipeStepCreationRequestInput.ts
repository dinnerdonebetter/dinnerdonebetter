/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepCreationRequestInput = {
  properties: {
    completionConditions: {
      type: 'array',
      contains: {
        type: 'RecipeStepCompletionConditionCreationRequestInput',
      },
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
    ingredients: {
      type: 'array',
      contains: {
        type: 'RecipeStepIngredientCreationRequestInput',
      },
    },
    instruments: {
      type: 'array',
      contains: {
        type: 'RecipeStepInstrumentCreationRequestInput',
      },
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
    preparationID: {
      type: 'string',
    },
    products: {
      type: 'array',
      contains: {
        type: 'RecipeStepProductCreationRequestInput',
      },
    },
    startTimerAutomatically: {
      type: 'boolean',
    },
    vessels: {
      type: 'array',
      contains: {
        type: 'RecipeStepVesselCreationRequestInput',
      },
    },
  },
} as const;
