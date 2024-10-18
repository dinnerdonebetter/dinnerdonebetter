// GENERATED CODE, DO NOT EDIT MANUALLY

import { User } from './User.gen';

export interface IHouseholdUserMembershipWithUser {
  archivedAt: string;
  belongsToHousehold: string;
  belongsToUser: User;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt: string;
}

export class HouseholdUserMembershipWithUser implements IHouseholdUserMembershipWithUser {
  archivedAt: string;
  belongsToHousehold: string;
  belongsToUser: User;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt: string;
  constructor(input: Partial<HouseholdUserMembershipWithUser> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.belongsToUser = input.belongsToUser || new User();
    this.createdAt = input.createdAt || '';
    this.defaultHousehold = input.defaultHousehold || false;
    this.householdRole = input.householdRole || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
  }
}
