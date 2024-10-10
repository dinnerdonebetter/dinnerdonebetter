// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdCreationRequestInput {
  addressLine1: string;
  addressLine2: string;
  city: string;
  contactPhone: string;
  country: string;
  latitude: number;
  longitude: number;
  name: string;
  state: string;
  zipCode: string;
}

export class HouseholdCreationRequestInput implements IHouseholdCreationRequestInput {
  addressLine1: string;
  addressLine2: string;
  city: string;
  contactPhone: string;
  country: string;
  latitude: number;
  longitude: number;
  name: string;
  state: string;
  zipCode: string;
  constructor(input: Partial<HouseholdCreationRequestInput> = {}) {
    this.addressLine1 = input.addressLine1 || '';
    this.addressLine2 = input.addressLine2 || '';
    this.city = input.city || '';
    this.contactPhone = input.contactPhone || '';
    this.country = input.country || '';
    this.latitude = input.latitude || 0;
    this.longitude = input.longitude || 0;
    this.name = input.name || '';
    this.state = input.state || '';
    this.zipCode = input.zipCode || '';
  }
}
