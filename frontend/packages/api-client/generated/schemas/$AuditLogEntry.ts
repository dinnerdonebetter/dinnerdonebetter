/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $AuditLogEntry = {
  properties: {
    belongsToHousehold: {
      type: 'string',
    },
    belongsToUser: {
      type: 'string',
    },
    changes: {
      type: 'ChangeLog',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    eventType: {
      type: 'string',
    },
    id: {
      type: 'string',
    },
    relevantID: {
      type: 'string',
    },
    resourceType: {
      type: 'string',
    },
  },
} as const;
