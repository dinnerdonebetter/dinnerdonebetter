// GENERATED CODE, DO NOT EDIT MANUALLY

import { HouseholdUserMembershipWithUser } from './HouseholdUserMembershipWithUser';

export interface IHousehold {
  createdAt: string;
  longitude?: number;
  members: HouseholdUserMembershipWithUser;
  name: string;
  addressLine1: string;
  contactPhone: string;
  country: string;
  id: string;
  lastUpdatedAt?: string;
  paymentProcessorCustomer: string;
  state: string;
  belongsToUser: string;
  city: string;
  billingStatus: string;
  subscriptionPlanID?: string;
  zipCode: string;
  addressLine2: string;
  archivedAt?: string;
  latitude?: number;
}

export class Household implements IHousehold {
  createdAt: string;
  longitude?: number;
  members: HouseholdUserMembershipWithUser;
  name: string;
  addressLine1: string;
  contactPhone: string;
  country: string;
  id: string;
  lastUpdatedAt?: string;
  paymentProcessorCustomer: string;
  state: string;
  belongsToUser: string;
  city: string;
  billingStatus: string;
  subscriptionPlanID?: string;
  zipCode: string;
  addressLine2: string;
  archivedAt?: string;
  latitude?: number;
  constructor(input: Partial<Household> = {}) {
    this.createdAt = input.createdAt = '';
    this.longitude = input.longitude;
    this.members = input.members = new HouseholdUserMembershipWithUser();
    this.name = input.name = '';
    this.addressLine1 = input.addressLine1 = '';
    this.contactPhone = input.contactPhone = '';
    this.country = input.country = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.paymentProcessorCustomer = input.paymentProcessorCustomer = '';
    this.state = input.state = '';
    this.belongsToUser = input.belongsToUser = '';
    this.city = input.city = '';
    this.billingStatus = input.billingStatus = '';
    this.subscriptionPlanID = input.subscriptionPlanID;
    this.zipCode = input.zipCode = '';
    this.addressLine2 = input.addressLine2 = '';
    this.archivedAt = input.archivedAt;
    this.latitude = input.latitude;
  }
}
