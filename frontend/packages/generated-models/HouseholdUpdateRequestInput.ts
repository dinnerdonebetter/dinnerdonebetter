// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdUpdateRequestInput {
  longitude?: number;
  name?: string;
  state?: string;
  addressLine1?: string;
  city?: string;
  latitude?: number;
  zipCode?: string;
  addressLine2?: string;
  contactPhone?: string;
  country?: string;
}

export class HouseholdUpdateRequestInput implements IHouseholdUpdateRequestInput {
  longitude?: number;
  name?: string;
  state?: string;
  addressLine1?: string;
  city?: string;
  latitude?: number;
  zipCode?: string;
  addressLine2?: string;
  contactPhone?: string;
  country?: string;
  constructor(input: Partial<HouseholdUpdateRequestInput> = {}) {
    this.longitude = input.longitude;
    this.name = input.name;
    this.state = input.state;
    this.addressLine1 = input.addressLine1;
    this.city = input.city;
    this.latitude = input.latitude;
    this.zipCode = input.zipCode;
    this.addressLine2 = input.addressLine2;
    this.contactPhone = input.contactPhone;
    this.country = input.country;
  }
}
