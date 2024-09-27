/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $PasswordUpdateInput = {
  properties: {
    currentPassword: {
      type: 'string',
      format: 'password',
    },
    newPassword: {
      type: 'string',
      format: 'password',
    },
    totpToken: {
      type: 'string',
    },
  },
} as const;
