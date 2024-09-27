/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $UsernameUpdateInput = {
  properties: {
    currentPassword: {
      type: 'string',
      format: 'password',
    },
    newUsername: {
      type: 'string',
    },
    totpToken: {
      type: 'string',
    },
  },
} as const;
