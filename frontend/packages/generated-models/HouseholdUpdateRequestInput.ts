// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdUpdateRequestInput {
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  contactPhone?: string;
  latitude?: number;
  name?: string;
  country?: string;
  longitude?: number;
  state?: string;
  zipCode?: string;
}

export class HouseholdUpdateRequestInput implements IHouseholdUpdateRequestInput {
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  contactPhone?: string;
  latitude?: number;
  name?: string;
  country?: string;
  longitude?: number;
  state?: string;
  zipCode?: string;
  constructor(input: Partial<HouseholdUpdateRequestInput> = {}) {
    this.addressLine1 = input.addressLine1;
    this.addressLine2 = input.addressLine2;
    this.city = input.city;
    this.contactPhone = input.contactPhone;
    this.latitude = input.latitude;
    this.name = input.name;
    this.country = input.country;
    this.longitude = input.longitude;
    this.state = input.state;
    this.zipCode = input.zipCode;
  }
}
