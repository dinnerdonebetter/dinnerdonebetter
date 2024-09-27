/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealPlanCreationRequestInput = {
  properties: {
    electionMethod: {
      type: 'string',
    },
    events: {
      type: 'array',
      contains: {
        type: 'MealPlanEventCreationRequestInput',
      },
    },
    notes: {
      type: 'string',
    },
    votingDeadline: {
      type: 'string',
      format: 'date-time',
    },
  },
} as const;
