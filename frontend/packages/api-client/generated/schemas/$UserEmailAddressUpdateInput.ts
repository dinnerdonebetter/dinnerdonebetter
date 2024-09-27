/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $UserEmailAddressUpdateInput = {
  properties: {
    currentPassword: {
      type: 'string',
      format: 'password',
    },
    newEmailAddress: {
      type: 'string',
    },
    totpToken: {
      type: 'string',
    },
  },
} as const;
