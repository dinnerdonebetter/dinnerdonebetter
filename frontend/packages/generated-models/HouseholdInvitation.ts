// GENERATED CODE, DO NOT EDIT MANUALLY

 import { Household } from './Household';
 import { User } from './User';


export interface IHouseholdInvitation {
   archivedAt?: string;
 lastUpdatedAt?: string;
 status: string;
 toName: string;
 toUser?: string;
 token: string;
 note: string;
 statusNote: string;
 toEmail: string;
 fromUser: User;
 id: string;
 createdAt: string;
 destinationHousehold: Household;
 expiresAt: string;

}

export class HouseholdInvitation implements IHouseholdInvitation {
   archivedAt?: string;
 lastUpdatedAt?: string;
 status: string;
 toName: string;
 toUser?: string;
 token: string;
 note: string;
 statusNote: string;
 toEmail: string;
 fromUser: User;
 id: string;
 createdAt: string;
 destinationHousehold: Household;
 expiresAt: string;
constructor(input: Partial<HouseholdInvitation> = {}) {
	 this.archivedAt = input.archivedAt;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.status = input.status = '';
 this.toName = input.toName = '';
 this.toUser = input.toUser;
 this.token = input.token = '';
 this.note = input.note = '';
 this.statusNote = input.statusNote = '';
 this.toEmail = input.toEmail = '';
 this.fromUser = input.fromUser = new User();
 this.id = input.id = '';
 this.createdAt = input.createdAt = '';
 this.destinationHousehold = input.destinationHousehold = new Household();
 this.expiresAt = input.expiresAt = '';
}
}