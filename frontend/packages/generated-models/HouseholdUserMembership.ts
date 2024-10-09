// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdUserMembership {
  belongsToUser: string;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  belongsToHousehold: string;
}

export class HouseholdUserMembership implements IHouseholdUserMembership {
  belongsToUser: string;
  createdAt: string;
  defaultHousehold: boolean;
  householdRole: string;
  id: string;
  lastUpdatedAt?: string;
  archivedAt?: string;
  belongsToHousehold: string;
  constructor(input: Partial<HouseholdUserMembership> = {}) {
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
    this.defaultHousehold = input.defaultHousehold = false;
    this.householdRole = input.householdRole = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
  }
}
