/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ServiceSetting = {
  properties: {
    adminsOnly: {
      type: 'boolean',
    },
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    defaultValue: {
      type: 'string',
    },
    description: {
      type: 'string',
    },
    enumeration: {
      type: 'array',
      contains: {
        type: 'string',
      },
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    name: {
      type: 'string',
    },
    type: {
      type: 'string',
    },
  },
} as const;
