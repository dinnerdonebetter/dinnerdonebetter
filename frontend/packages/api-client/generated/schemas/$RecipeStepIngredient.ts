/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepIngredient = {
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
    ingredient: {
      type: 'ValidIngredient',
    },
    ingredientNotes: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    maximumQuantity: {
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
    name: {
      type: 'string',
    },
    optionIndex: {
      type: 'number',
      format: 'int64',
    },
    optional: {
      type: 'boolean',
    },
    productOfRecipeID: {
      type: 'string',
    },
    productPercentageToUse: {
      type: 'number',
      format: 'double',
    },
    quantityNotes: {
      type: 'string',
    },
    recipeStepProductID: {
      type: 'string',
    },
    toTaste: {
      type: 'boolean',
    },
    vesselIndex: {
      type: 'number',
      format: 'int64',
    },
  },
} as const;
