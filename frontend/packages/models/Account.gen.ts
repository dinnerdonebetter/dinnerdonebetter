// GENERATED CODE, DO NOT EDIT MANUALLY

import { AccountUserMembershipWithUser } from './AccountUserMembershipWithUser.gen';

export interface IAccount {
  addressLine1: string;
  addressLine2: string;
  archivedAt: string;
  belongsToUser: string;
  billingStatus: string;
  city: string;
  contactPhone: string;
  country: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  latitude: number;
  longitude: number;
  members: AccountUserMembershipWithUser[];
  name: string;
  paymentProcessorCustomer: string;
  state: string;
  subscriptionPlanID: string;
  zipCode: string;
}

export class Account implements IAccount {
  addressLine1: string;
  addressLine2: string;
  archivedAt: string;
  belongsToUser: string;
  billingStatus: string;
  city: string;
  contactPhone: string;
  country: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  latitude: number;
  longitude: number;
  members: AccountUserMembershipWithUser[];
  name: string;
  paymentProcessorCustomer: string;
  state: string;
  subscriptionPlanID: string;
  zipCode: string;
  constructor(input: Partial<Account> = {}) {
    this.addressLine1 = input.addressLine1 || '';
    this.addressLine2 = input.addressLine2 || '';
    this.archivedAt = input.archivedAt || '';
    this.belongsToUser = input.belongsToUser || '';
    this.billingStatus = input.billingStatus || '';
    this.city = input.city || '';
    this.contactPhone = input.contactPhone || '';
    this.country = input.country || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.latitude = input.latitude || 0;
    this.longitude = input.longitude || 0;
    this.members = input.members || [];
    this.name = input.name || '';
    this.paymentProcessorCustomer = input.paymentProcessorCustomer || '';
    this.state = input.state || '';
    this.subscriptionPlanID = input.subscriptionPlanID || '';
    this.zipCode = input.zipCode || '';
  }
}
