// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IHouseholdCreationRequestInput {
   zipCode: string;
 addressLine1: string;
 contactPhone: string;
 country: string;
 latitude?: number;
 name: string;
 addressLine2: string;
 city: string;
 longitude?: number;
 state: string;

}

export class HouseholdCreationRequestInput implements IHouseholdCreationRequestInput {
   zipCode: string;
 addressLine1: string;
 contactPhone: string;
 country: string;
 latitude?: number;
 name: string;
 addressLine2: string;
 city: string;
 longitude?: number;
 state: string;
constructor(input: Partial<HouseholdCreationRequestInput> = {}) {
	 this.zipCode = input.zipCode = '';
 this.addressLine1 = input.addressLine1 = '';
 this.contactPhone = input.contactPhone = '';
 this.country = input.country = '';
 this.latitude = input.latitude;
 this.name = input.name = '';
 this.addressLine2 = input.addressLine2 = '';
 this.city = input.city = '';
 this.longitude = input.longitude;
 this.state = input.state = '';
}
}