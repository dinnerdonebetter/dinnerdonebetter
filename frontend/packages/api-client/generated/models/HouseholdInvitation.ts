/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Household } from './Household';
import type { User } from './User';
export type HouseholdInvitation = {
  archivedAt?: string;
  createdAt?: string;
  destinationHousehold?: Household;
  expiresAt?: string;
  fromUser?: User;
  id?: string;
  lastUpdatedAt?: string;
  note?: string;
  status?: string;
  statusNote?: string;
  toEmail?: string;
  toName?: string;
  toUser?: string;
  token?: string;
};
