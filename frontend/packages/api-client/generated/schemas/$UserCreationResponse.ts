/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $UserCreationResponse = {
  properties: {
    accountStatus: {
      type: 'string',
    },
    avatar: {
      type: 'string',
    },
    birthday: {
      type: 'string',
      format: 'date-time',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    createdUserID: {
      type: 'string',
    },
    emailAddress: {
      type: 'string',
      format: 'email',
    },
    firstName: {
      type: 'string',
    },
    isAdmin: {
      type: 'boolean',
    },
    lastName: {
      type: 'string',
    },
    qrCode: {
      type: 'string',
    },
    twoFactorSecret: {
      type: 'string',
    },
    username: {
      type: 'string',
    },
  },
} as const;
