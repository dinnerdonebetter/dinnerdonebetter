/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipePrepTask = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToRecipe: {
      type: 'string',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    description: {
      type: 'string',
    },
    explicitStorageInstructions: {
      type: 'string',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    maximumStorageTemperatureInCelsius: {
      type: 'number',
      format: 'double',
    },
    maximumTimeBufferBeforeRecipeInSeconds: {
      type: 'number',
      format: 'int64',
    },
    minimumStorageTemperatureInCelsius: {
      type: 'number',
      format: 'double',
    },
    minimumTimeBufferBeforeRecipeInSeconds: {
      type: 'number',
      format: 'int64',
    },
    name: {
      type: 'string',
    },
    notes: {
      type: 'string',
    },
    optional: {
      type: 'boolean',
    },
    recipeSteps: {
      type: 'array',
      contains: {
        type: 'RecipePrepTaskStep',
      },
    },
    storageType: {
      type: 'string',
    },
  },
} as const;
