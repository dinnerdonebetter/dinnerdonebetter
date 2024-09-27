/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipePrepTaskCreationRequestInput = {
  properties: {
    belongsToRecipe: {
      type: 'string',
    },
    description: {
      type: 'string',
    },
    explicitStorageInstructions: {
      type: 'string',
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
        type: 'RecipePrepTaskStepCreationRequestInput',
      },
    },
    storageType: {
      type: 'string',
    },
  },
} as const;
