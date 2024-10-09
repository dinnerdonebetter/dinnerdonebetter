// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdCreationRequestInput {
  addressLine2: string;
  city: string;
  latitude?: number;
  longitude?: number;
  state: string;
  addressLine1: string;
  country: string;
  name: string;
  zipCode: string;
  contactPhone: string;
}

export class HouseholdCreationRequestInput implements IHouseholdCreationRequestInput {
  addressLine2: string;
  city: string;
  latitude?: number;
  longitude?: number;
  state: string;
  addressLine1: string;
  country: string;
  name: string;
  zipCode: string;
  contactPhone: string;
  constructor(input: Partial<HouseholdCreationRequestInput> = {}) {
    this.addressLine2 = input.addressLine2 = '';
    this.city = input.city = '';
    this.latitude = input.latitude;
    this.longitude = input.longitude;
    this.state = input.state = '';
    this.addressLine1 = input.addressLine1 = '';
    this.country = input.country = '';
    this.name = input.name = '';
    this.zipCode = input.zipCode = '';
    this.contactPhone = input.contactPhone = '';
  }
}
