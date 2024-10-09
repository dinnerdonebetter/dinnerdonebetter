// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdUpdateRequestInput {
  name: string;
  state: string;
  zipCode: string;
  addressLine2: string;
  longitude: number;
  contactPhone: string;
  country: string;
  latitude: number;
  addressLine1: string;
  city: string;
}

export class HouseholdUpdateRequestInput implements IHouseholdUpdateRequestInput {
  name: string;
  state: string;
  zipCode: string;
  addressLine2: string;
  longitude: number;
  contactPhone: string;
  country: string;
  latitude: number;
  addressLine1: string;
  city: string;
  constructor(input: Partial<HouseholdUpdateRequestInput> = {}) {
    this.name = input.name || '';
    this.state = input.state || '';
    this.zipCode = input.zipCode || '';
    this.addressLine2 = input.addressLine2 || '';
    this.longitude = input.longitude || 0;
    this.contactPhone = input.contactPhone || '';
    this.country = input.country || '';
    this.latitude = input.latitude || 0;
    this.addressLine1 = input.addressLine1 || '';
    this.city = input.city || '';
  }
}
