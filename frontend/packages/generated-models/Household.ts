// GENERATED CODE, DO NOT EDIT MANUALLY

 import { HouseholdUserMembershipWithUser } from './HouseholdUserMembershipWithUser';


export interface IHousehold {
   longitude?: number;
 name: string;
 subscriptionPlanID?: string;
 zipCode: string;
 contactPhone: string;
 id: string;
 country: string;
 createdAt: string;
 members: HouseholdUserMembershipWithUser;
 state: string;
 addressLine2: string;
 billingStatus: string;
 belongsToUser: string;
 city: string;
 latitude?: number;
 addressLine1: string;
 archivedAt?: string;
 lastUpdatedAt?: string;
 paymentProcessorCustomer: string;

}

export class Household implements IHousehold {
   longitude?: number;
 name: string;
 subscriptionPlanID?: string;
 zipCode: string;
 contactPhone: string;
 id: string;
 country: string;
 createdAt: string;
 members: HouseholdUserMembershipWithUser;
 state: string;
 addressLine2: string;
 billingStatus: string;
 belongsToUser: string;
 city: string;
 latitude?: number;
 addressLine1: string;
 archivedAt?: string;
 lastUpdatedAt?: string;
 paymentProcessorCustomer: string;
constructor(input: Partial<Household> = {}) {
	 this.longitude = input.longitude;
 this.name = input.name = '';
 this.subscriptionPlanID = input.subscriptionPlanID;
 this.zipCode = input.zipCode = '';
 this.contactPhone = input.contactPhone = '';
 this.id = input.id = '';
 this.country = input.country = '';
 this.createdAt = input.createdAt = '';
 this.members = input.members = new HouseholdUserMembershipWithUser();
 this.state = input.state = '';
 this.addressLine2 = input.addressLine2 = '';
 this.billingStatus = input.billingStatus = '';
 this.belongsToUser = input.belongsToUser = '';
 this.city = input.city = '';
 this.latitude = input.latitude;
 this.addressLine1 = input.addressLine1 = '';
 this.archivedAt = input.archivedAt;
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.paymentProcessorCustomer = input.paymentProcessorCustomer = '';
}
}