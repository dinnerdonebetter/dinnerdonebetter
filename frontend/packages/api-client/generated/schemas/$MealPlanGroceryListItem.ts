/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealPlanGroceryListItem = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToMealPlan: {
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
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    maximumQuantityNeeded: {
      type: 'number',
      format: 'double',
    },
    measurementUnit: {
      type: 'ValidMeasurementUnit',
    },
    minimumQuantityNeeded: {
      type: 'number',
      format: 'double',
    },
    purchasePrice: {
      type: 'number',
      format: 'double',
    },
    purchasedMeasurementUnit: {
      type: 'ValidMeasurementUnit',
    },
    purchasedUPC: {
      type: 'string',
    },
    quantityPurchased: {
      type: 'number',
      format: 'double',
    },
    status: {
      type: 'string',
    },
    statusExplanation: {
      type: 'string',
    },
  },
} as const;
