// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdCreationRequestInput {
  addressLine2: string;
  city: string;
  contactPhone: string;
  latitude: number;
  name: string;
  state: string;
  addressLine1: string;
  country: string;
  longitude: number;
  zipCode: string;
}

export class HouseholdCreationRequestInput implements IHouseholdCreationRequestInput {
  addressLine2: string;
  city: string;
  contactPhone: string;
  latitude: number;
  name: string;
  state: string;
  addressLine1: string;
  country: string;
  longitude: number;
  zipCode: string;
  constructor(input: Partial<HouseholdCreationRequestInput> = {}) {
    this.addressLine2 = input.addressLine2 || '';
    this.city = input.city || '';
    this.contactPhone = input.contactPhone || '';
    this.latitude = input.latitude || 0;
    this.name = input.name || '';
    this.state = input.state || '';
    this.addressLine1 = input.addressLine1 || '';
    this.country = input.country || '';
    this.longitude = input.longitude || 0;
    this.zipCode = input.zipCode || '';
  }
}
