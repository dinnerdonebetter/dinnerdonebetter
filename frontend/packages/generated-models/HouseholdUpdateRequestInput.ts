// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IHouseholdUpdateRequestInput {
   addressLine2?: string;
 contactPhone?: string;
 country?: string;
 latitude?: number;
 longitude?: number;
 addressLine1?: string;
 city?: string;
 name?: string;
 state?: string;
 zipCode?: string;

}

export class HouseholdUpdateRequestInput implements IHouseholdUpdateRequestInput {
   addressLine2?: string;
 contactPhone?: string;
 country?: string;
 latitude?: number;
 longitude?: number;
 addressLine1?: string;
 city?: string;
 name?: string;
 state?: string;
 zipCode?: string;
constructor(input: Partial<HouseholdUpdateRequestInput> = {}) {
	 this.addressLine2 = input.addressLine2;
 this.contactPhone = input.contactPhone;
 this.country = input.country;
 this.latitude = input.latitude;
 this.longitude = input.longitude;
 this.addressLine1 = input.addressLine1;
 this.city = input.city;
 this.name = input.name;
 this.state = input.state;
 this.zipCode = input.zipCode;
}
}