// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IHouseholdUserMembership {
   archivedAt: string;
 belongsToHousehold: string;
 belongsToUser: string;
 createdAt: string;
 defaultHousehold: boolean;
 householdRole: string;
 id: string;
 lastUpdatedAt: string;

}

export class HouseholdUserMembership implements IHouseholdUserMembership {
   archivedAt: string;
 belongsToHousehold: string;
 belongsToUser: string;
 createdAt: string;
 defaultHousehold: boolean;
 householdRole: string;
 id: string;
 lastUpdatedAt: string;
constructor(input: Partial<HouseholdUserMembership> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.belongsToHousehold = input.belongsToHousehold || '';
 this.belongsToUser = input.belongsToUser || '';
 this.createdAt = input.createdAt || '';
 this.defaultHousehold = input.defaultHousehold || false;
 this.householdRole = input.householdRole || '';
 this.id = input.id || '';
 this.lastUpdatedAt = input.lastUpdatedAt || '';
}
}