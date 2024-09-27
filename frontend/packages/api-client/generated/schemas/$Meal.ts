/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $Meal = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    components: {
      type: 'array',
      contains: {
        type: 'MealComponent',
      },
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    createdByUser: {
      type: 'string',
    },
    description: {
      type: 'string',
    },
    eligibleForMealPlans: {
      type: 'boolean',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
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
