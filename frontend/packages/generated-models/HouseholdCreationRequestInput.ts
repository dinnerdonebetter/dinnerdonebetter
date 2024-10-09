// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdCreationRequestInput {
  city: string;
  country: string;
  latitude?: number;
  longitude?: number;
  name: string;
  state: string;
  addressLine1: string;
  addressLine2: string;
  contactPhone: string;
  zipCode: string;
}

export class HouseholdCreationRequestInput implements IHouseholdCreationRequestInput {
  city: string;
  country: string;
  latitude?: number;
  longitude?: number;
  name: string;
  state: string;
  addressLine1: string;
  addressLine2: string;
  contactPhone: string;
  zipCode: string;
  constructor(input: Partial<HouseholdCreationRequestInput> = {}) {
    this.city = input.city = '';
    this.country = input.country = '';
    this.latitude = input.latitude;
    this.longitude = input.longitude;
    this.name = input.name = '';
    this.state = input.state = '';
    this.addressLine1 = input.addressLine1 = '';
    this.addressLine2 = input.addressLine2 = '';
    this.contactPhone = input.contactPhone = '';
    this.zipCode = input.zipCode = '';
  }
}
