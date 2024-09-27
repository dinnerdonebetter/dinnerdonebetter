/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealCreationRequestInput = {
  properties: {
    components: {
      type: 'array',
      contains: {
        type: 'MealComponentCreationRequestInput',
      },
    },
    description: {
      type: 'string',
    },
    eligibleForMealPlans: {
      type: 'boolean',
    },
    maximumEstimatedPortions: {
      type: 'number',
      format: 'double',
    },
    minimumEstimatedPortions: {
      type: 'number',
      format: 'double',
    },
    name: {
      type: 'string',
    },
  },
} as const;
