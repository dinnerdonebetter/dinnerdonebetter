// GENERATED CODE, DO NOT EDIT MANUALLY

import { User } from './User';

export interface IHouseholdUserMembershipWithUser {
  belongsToHousehold: string;
  belongsToUser?: User;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
}

export class HouseholdUserMembershipWithUser implements IHouseholdUserMembershipWithUser {
  belongsToHousehold: string;
  belongsToUser?: User;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  constructor(input: Partial<HouseholdUserMembershipWithUser> = {}) {
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.belongsToUser = input.belongsToUser;
    this.createdAt = input.createdAt = '';
    this.defaultHousehold = input.defaultHousehold = false;
    this.householdRole = input.householdRole = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
  }
}
