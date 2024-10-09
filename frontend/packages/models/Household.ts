// GENERATED CODE, DO NOT EDIT MANUALLY

import { HouseholdUserMembershipWithUser } from './HouseholdUserMembershipWithUser';

export interface IHousehold {
  zipCode: string;
  addressLine1: string;
  members: HouseholdUserMembershipWithUser[];
  paymentProcessorCustomer: string;
  city: string;
  subscriptionPlanID: string;
  contactPhone: string;
  country: string;
  createdAt: string;
  lastUpdatedAt: string;
  addressLine2: string;
  archivedAt: string;
  belongsToUser: string;
  longitude: number;
  name: string;
  state: string;
  billingStatus: string;
  id: string;
  latitude: number;
}

export class Household implements IHousehold {
  zipCode: string;
  addressLine1: string;
  members: HouseholdUserMembershipWithUser[];
  paymentProcessorCustomer: string;
  city: string;
  subscriptionPlanID: string;
  contactPhone: string;
  country: string;
  createdAt: string;
  lastUpdatedAt: string;
  addressLine2: string;
  archivedAt: string;
  belongsToUser: string;
  longitude: number;
  name: string;
  state: string;
  billingStatus: string;
  id: string;
  latitude: number;
  constructor(input: Partial<Household> = {}) {
    this.zipCode = input.zipCode || '';
    this.addressLine1 = input.addressLine1 || '';
    this.members = input.members || [];
    this.paymentProcessorCustomer = input.paymentProcessorCustomer || '';
    this.city = input.city || '';
    this.subscriptionPlanID = input.subscriptionPlanID || '';
    this.contactPhone = input.contactPhone || '';
    this.country = input.country || '';
    this.createdAt = input.createdAt || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.addressLine2 = input.addressLine2 || '';
    this.archivedAt = input.archivedAt || '';
    this.belongsToUser = input.belongsToUser || '';
    this.longitude = input.longitude || 0;
    this.name = input.name || '';
    this.state = input.state || '';
    this.billingStatus = input.billingStatus || '';
    this.id = input.id || '';
    this.latitude = input.latitude || 0;
  }
}
