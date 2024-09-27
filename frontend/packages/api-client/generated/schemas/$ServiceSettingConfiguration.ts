/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ServiceSettingConfiguration = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToHousehold: {
      type: 'string',
    },
    belongsToUser: {
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
    notes: {
      type: 'string',
    },
    serviceSetting: {
      type: 'ServiceSetting',
    },
    value: {
      type: 'string',
    },
  },
} as const;
