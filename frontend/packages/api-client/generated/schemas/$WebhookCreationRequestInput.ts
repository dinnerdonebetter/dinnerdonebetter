/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $WebhookCreationRequestInput = {
  properties: {
    contentType: {
      type: 'string',
    },
    events: {
      type: 'array',
      contains: {
        type: 'string',
      },
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
