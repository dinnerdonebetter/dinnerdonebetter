// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdUpdateRequestInput {
  name?: string;
  state?: string;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  contactPhone?: string;
  country?: string;
  latitude?: number;
  longitude?: number;
  zipCode?: string;
}

export class HouseholdUpdateRequestInput implements IHouseholdUpdateRequestInput {
  name?: string;
  state?: string;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  contactPhone?: string;
  country?: string;
  latitude?: number;
  longitude?: number;
  zipCode?: string;
  constructor(input: Partial<HouseholdUpdateRequestInput> = {}) {
    this.name = input.name;
    this.state = input.state;
    this.addressLine1 = input.addressLine1;
    this.addressLine2 = input.addressLine2;
    this.city = input.city;
    this.contactPhone = input.contactPhone;
    this.country = input.country;
    this.latitude = input.latitude;
    this.longitude = input.longitude;
    this.zipCode = input.zipCode;
  }
}
