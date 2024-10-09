// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdUserMembership {
  belongsToHousehold: string;
  belongsToUser: string;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
}

export class HouseholdUserMembership implements IHouseholdUserMembership {
  belongsToHousehold: string;
  belongsToUser: string;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  constructor(input: Partial<HouseholdUserMembership> = {}) {
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
    this.defaultHousehold = input.defaultHousehold = false;
    this.householdRole = input.householdRole = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
  }
}
