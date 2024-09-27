/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealPlanEvent = {
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
    endsAt: {
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
    mealName: {
      type: 'string',
    },
    notes: {
      type: 'string',
    },
    options: {
      type: 'array',
      contains: {
        type: 'MealPlanOption',
      },
    },
    startsAt: {
      type: 'string',
      format: 'date-time',
    },
  },
} as const;
