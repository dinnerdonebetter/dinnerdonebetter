/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStep = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToRecipe: {
      type: 'string',
    },
    completionConditions: {
      type: 'array',
      contains: {
        type: 'RecipeStepCompletionCondition',
      },
    },
    conditionExpression: {
      type: 'string',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    explicitInstructions: {
      type: 'string',
    },
    id: {
      type: 'string',
    },
    index: {
      type: 'number',
      format: 'int64',
    },
    ingredients: {
      type: 'array',
      contains: {
        type: 'RecipeStepIngredient',
      },
    },
    instruments: {
      type: 'array',
      contains: {
        type: 'RecipeStepInstrument',
      },
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    maximumEstimatedTimeInSeconds: {
      type: 'number',
      format: 'int64',
    },
    maximumTemperatureInCelsius: {
      type: 'number',
      format: 'double',
    },
    media: {
      type: 'array',
      contains: {
        type: 'RecipeMedia',
      },
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
    products: {
      type: 'array',
      contains: {
        type: 'RecipeStepProduct',
      },
    },
    startTimerAutomatically: {
      type: 'boolean',
    },
    vessels: {
      type: 'array',
      contains: {
        type: 'RecipeStepVessel',
      },
    },
  },
} as const;
