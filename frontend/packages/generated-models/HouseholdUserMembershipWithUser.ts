// GENERATED CODE, DO NOT EDIT MANUALLY

import { User } from './User';

export interface IHouseholdUserMembershipWithUser {
  lastUpdatedAt?: string;
  archivedAt?: string;
  belongsToHousehold: string;
  belongsToUser?: User;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
}

export class HouseholdUserMembershipWithUser implements IHouseholdUserMembershipWithUser {
  lastUpdatedAt?: string;
  archivedAt?: string;
  belongsToHousehold: string;
  belongsToUser?: User;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  constructor(input: Partial<HouseholdUserMembershipWithUser> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.belongsToUser = input.belongsToUser;
    this.createdAt = input.createdAt = '';
    this.defaultHousehold = input.defaultHousehold = false;
    this.householdRole = input.householdRole = '';
    this.id = input.id = '';
  }
}
