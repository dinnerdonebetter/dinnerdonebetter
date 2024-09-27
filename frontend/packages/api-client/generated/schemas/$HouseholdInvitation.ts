/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $HouseholdInvitation = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    destinationHousehold: {
      type: 'Household',
    },
    expiresAt: {
      type: 'string',
      format: 'date-time',
    },
    fromUser: {
      type: 'User',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    note: {
      type: 'string',
    },
    status: {
      type: 'string',
    },
    statusNote: {
      type: 'string',
    },
    toEmail: {
      type: 'string',
    },
    toName: {
      type: 'string',
    },
    toUser: {
      type: 'string',
    },
    token: {
      type: 'string',
    },
  },
} as const;
