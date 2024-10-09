// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdCreationRequestInput {
  longitude?: number;
  state: string;
  contactPhone: string;
  country: string;
  latitude?: number;
  name: string;
  zipCode: string;
  addressLine1: string;
  addressLine2: string;
  city: string;
}

export class HouseholdCreationRequestInput implements IHouseholdCreationRequestInput {
  longitude?: number;
  state: string;
  contactPhone: string;
  country: string;
  latitude?: number;
  name: string;
  zipCode: string;
  addressLine1: string;
  addressLine2: string;
  city: string;
  constructor(input: Partial<HouseholdCreationRequestInput> = {}) {
    this.longitude = input.longitude;
    this.state = input.state = '';
    this.contactPhone = input.contactPhone = '';
    this.country = input.country = '';
    this.latitude = input.latitude;
    this.name = input.name = '';
    this.zipCode = input.zipCode = '';
    this.addressLine1 = input.addressLine1 = '';
    this.addressLine2 = input.addressLine2 = '';
    this.city = input.city = '';
  }
}
