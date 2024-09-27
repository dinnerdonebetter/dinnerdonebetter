/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepVessel = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToRecipeStep: {
      type: 'string',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    maximumQuantity: {
      type: 'number',
      format: 'int64',
    },
    minimumQuantity: {
      type: 'number',
      format: 'int64',
    },
    name: {
      type: 'string',
    },
    notes: {
      type: 'string',
    },
    recipeStepProductID: {
      type: 'string',
    },
    unavailableAfterStep: {
      type: 'boolean',
    },
    vessel: {
      type: 'ValidVessel',
    },
    vesselPreposition: {
      type: 'string',
    },
  },
} as const;
