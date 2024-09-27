/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $Household = {
  properties: {
    addressLine1: {
      type: 'string',
    },
    addressLine2: {
      type: 'string',
    },
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToUser: {
      type: 'string',
    },
    billingStatus: {
      type: 'string',
    },
    city: {
      type: 'string',
    },
    contactPhone: {
      type: 'string',
    },
    country: {
      type: 'string',
    },
    createdAt: {
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
    latitude: {
      type: 'number',
      format: 'double',
    },
    longitude: {
      type: 'number',
      format: 'double',
    },
    members: {
      type: 'array',
      contains: {
        type: 'HouseholdUserMembershipWithUser',
      },
    },
    name: {
      type: 'string',
    },
    paymentProcessorCustomer: {
      type: 'string',
    },
    state: {
      type: 'string',
    },
    subscriptionPlanID: {
      type: 'string',
    },
    zipCode: {
      type: 'string',
    },
  },
} as const;
