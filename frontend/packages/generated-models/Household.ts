// GENERATED CODE, DO NOT EDIT MANUALLY

import { HouseholdUserMembershipWithUser } from './HouseholdUserMembershipWithUser';

export interface IHousehold {
  addressLine1: string;
  addressLine2: string;
  city: string;
  subscriptionPlanID?: string;
  billingStatus: string;
  createdAt: string;
  id: string;
  state: string;
  longitude?: number;
  members: HouseholdUserMembershipWithUser;
  paymentProcessorCustomer: string;
  archivedAt?: string;
  belongsToUser: string;
  country: string;
  lastUpdatedAt?: string;
  latitude?: number;
  contactPhone: string;
  name: string;
  zipCode: string;
}

export class Household implements IHousehold {
  addressLine1: string;
  addressLine2: string;
  city: string;
  subscriptionPlanID?: string;
  billingStatus: string;
  createdAt: string;
  id: string;
  state: string;
  longitude?: number;
  members: HouseholdUserMembershipWithUser;
  paymentProcessorCustomer: string;
  archivedAt?: string;
  belongsToUser: string;
  country: string;
  lastUpdatedAt?: string;
  latitude?: number;
  contactPhone: string;
  name: string;
  zipCode: string;
  constructor(input: Partial<Household> = {}) {
    this.addressLine1 = input.addressLine1 = '';
    this.addressLine2 = input.addressLine2 = '';
    this.city = input.city = '';
    this.subscriptionPlanID = input.subscriptionPlanID;
    this.billingStatus = input.billingStatus = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.state = input.state = '';
    this.longitude = input.longitude;
    this.members = input.members = new HouseholdUserMembershipWithUser();
    this.paymentProcessorCustomer = input.paymentProcessorCustomer = '';
    this.archivedAt = input.archivedAt;
    this.belongsToUser = input.belongsToUser = '';
    this.country = input.country = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.latitude = input.latitude;
    this.contactPhone = input.contactPhone = '';
    this.name = input.name = '';
    this.zipCode = input.zipCode = '';
  }
}
