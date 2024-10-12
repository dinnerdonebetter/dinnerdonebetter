// GENERATED CODE, DO NOT EDIT MANUALLY

 import { Household } from './Household';
 import { User } from './User';


export interface IHouseholdInvitation {
   archivedAt: string;
 createdAt: string;
 destinationHousehold: Household;
 expiresAt: string;
 fromUser: User;
 id: string;
 lastUpdatedAt: string;
 note: string;
 status: string;
 statusNote: string;
 toEmail: string;
 toName: string;
 toUser: string;
 token: string;

}

export class HouseholdInvitation implements IHouseholdInvitation {
   archivedAt: string;
 createdAt: string;
 destinationHousehold: Household;
 expiresAt: string;
 fromUser: User;
 id: string;
 lastUpdatedAt: string;
 note: string;
 status: string;
 statusNote: string;
 toEmail: string;
 toName: string;
 toUser: string;
 token: string;
constructor(input: Partial<HouseholdInvitation> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.createdAt = input.createdAt || '';
 this.destinationHousehold = input.destinationHousehold || new Household();
 this.expiresAt = input.expiresAt || '';
 this.fromUser = input.fromUser || new User();
 this.id = input.id || '';
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.note = input.note || '';
 this.status = input.status || '';
 this.statusNote = input.statusNote || '';
 this.toEmail = input.toEmail || '';
 this.toName = input.toName || '';
 this.toUser = input.toUser || '';
 this.token = input.token || '';
}
}