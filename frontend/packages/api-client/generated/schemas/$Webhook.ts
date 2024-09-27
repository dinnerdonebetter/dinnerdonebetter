/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $Webhook = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToHousehold: {
      type: 'string',
    },
    contentType: {
      type: 'string',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    events: {
      type: 'array',
      contains: {
        type: 'WebhookTriggerEvent',
      },
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    method: {
      type: 'string',
    },
    name: {
      type: 'string',
    },
    url: {
      type: 'string',
      format: 'uri',
    },
  },
} as const;
