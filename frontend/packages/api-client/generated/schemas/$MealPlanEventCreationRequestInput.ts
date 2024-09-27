/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealPlanEventCreationRequestInput = {
  properties: {
    endsAt: {
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
        type: 'MealPlanOptionCreationRequestInput',
      },
    },
    startsAt: {
      type: 'string',
      format: 'date-time',
    },
  },
} as const;
