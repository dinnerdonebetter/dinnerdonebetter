// GENERATED CODE, DO NOT EDIT MANUALLY

import { HouseholdUserMembershipWithUser } from './HouseholdUserMembershipWithUser';

export interface IHousehold {
  id: string;
  name: string;
  subscriptionPlanID?: string;
  belongsToUser: string;
  createdAt: string;
  lastUpdatedAt?: string;
  latitude?: number;
  state: string;
  addressLine2: string;
  billingStatus: string;
  members: HouseholdUserMembershipWithUser;
  paymentProcessorCustomer: string;
  zipCode: string;
  addressLine1: string;
  archivedAt?: string;
  city: string;
  contactPhone: string;
  country: string;
  longitude?: number;
}

export class Household implements IHousehold {
  id: string;
  name: string;
  subscriptionPlanID?: string;
  belongsToUser: string;
  createdAt: string;
  lastUpdatedAt?: string;
  latitude?: number;
  state: string;
  addressLine2: string;
  billingStatus: string;
  members: HouseholdUserMembershipWithUser;
  paymentProcessorCustomer: string;
  zipCode: string;
  addressLine1: string;
  archivedAt?: string;
  city: string;
  contactPhone: string;
  country: string;
  longitude?: number;
  constructor(input: Partial<Household> = {}) {
    this.id = input.id = '';
    this.name = input.name = '';
    this.subscriptionPlanID = input.subscriptionPlanID;
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.latitude = input.latitude;
    this.state = input.state = '';
    this.addressLine2 = input.addressLine2 = '';
    this.billingStatus = input.billingStatus = '';
    this.members = input.members = new HouseholdUserMembershipWithUser();
    this.paymentProcessorCustomer = input.paymentProcessorCustomer = '';
    this.zipCode = input.zipCode = '';
    this.addressLine1 = input.addressLine1 = '';
    this.archivedAt = input.archivedAt;
    this.city = input.city = '';
    this.contactPhone = input.contactPhone = '';
    this.country = input.country = '';
    this.longitude = input.longitude;
  }
}
