// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdUserMembership {
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  belongsToHousehold: string;
  belongsToUser: string;
  createdAt: string;
}

export class HouseholdUserMembership implements IHouseholdUserMembership {
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  belongsToHousehold: string;
  belongsToUser: string;
  createdAt: string;
  constructor(input: Partial<HouseholdUserMembership> = {}) {
    this.defaultHousehold = input.defaultHousehold = false;
    this.householdRole = input.householdRole = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
  }
}
