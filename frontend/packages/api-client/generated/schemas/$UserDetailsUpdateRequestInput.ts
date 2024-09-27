/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $UserDetailsUpdateRequestInput = {
  properties: {
    birthday: {
      type: 'string',
      format: 'date-time',
    },
    currentPassword: {
      type: 'string',
      format: 'password',
    },
    firstName: {
      type: 'string',
    },
    lastName: {
      type: 'string',
    },
    totpToken: {
      type: 'string',
    },
  },
} as const;
