/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepInstrumentCreationRequestInput = {
  properties: {
    instrumentID: {
      type: 'string',
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
    productOfRecipeStepIndex: {
      type: 'number',
      format: 'int64',
    },
    productOfRecipeStepProductIndex: {
      type: 'number',
      format: 'int64',
    },
    recipeStepProductID: {
      type: 'string',
    },
  },
} as const;
