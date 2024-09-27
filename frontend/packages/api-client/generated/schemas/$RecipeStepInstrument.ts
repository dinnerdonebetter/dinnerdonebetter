/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepInstrument = {
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
    instrument: {
      type: 'ValidInstrument',
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
    optionIndex: {
      type: 'number',
      format: 'int64',
    },
    optional: {
      type: 'boolean',
    },
    preferenceRank: {
      type: 'number',
      format: 'int32',
    },
    recipeStepProductID: {
      type: 'string',
    },
  },
} as const;
