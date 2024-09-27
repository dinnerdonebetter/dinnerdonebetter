/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealPlanGroceryListItemCreationRequestInput = {
  properties: {
    belongsToMealPlan: {
      type: 'string',
    },
    maximumQuantityNeeded: {
      type: 'number',
      format: 'double',
    },
    minimumQuantityNeeded: {
      type: 'number',
      format: 'double',
    },
    purchasePrice: {
      type: 'number',
      format: 'double',
    },
    purchasedMeasurementUnitID: {
      type: 'string',
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
    validIngredientID: {
      type: 'string',
    },
    validMeasurementUnitID: {
      type: 'string',
    },
  },
} as const;
