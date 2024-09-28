/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { HouseholdUserMembershipWithUser } from './HouseholdUserMembershipWithUser';
export type Household = {
  addressLine1?: string;
  addressLine2?: string;
  archivedAt?: string;
  belongsToUser?: string;
  billingStatus?: string;
  city?: string;
  contactPhone?: string;
  country?: string;
  createdAt?: string;
  id?: string;
  lastUpdatedAt?: string;
  latitude?: number;
  longitude?: number;
  members?: Array<HouseholdUserMembershipWithUser>;
  name?: string;
  paymentProcessorCustomer?: string;
  state?: string;
  subscriptionPlanID?: string;
  zipCode?: string;
};
