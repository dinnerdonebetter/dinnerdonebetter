// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument';

export interface IHouseholdInstrumentOwnership {
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt?: string;
  notes: string;
  quantity: number;
  archivedAt?: string;
  belongsToHousehold: string;
  createdAt: string;
}

export class HouseholdInstrumentOwnership implements IHouseholdInstrumentOwnership {
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt?: string;
  notes: string;
  quantity: number;
  archivedAt?: string;
  belongsToHousehold: string;
  createdAt: string;
  constructor(input: Partial<HouseholdInstrumentOwnership> = {}) {
    this.id = input.id = '';
    this.instrument = input.instrument = new ValidInstrument();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.quantity = input.quantity = 0;
    this.archivedAt = input.archivedAt;
    this.belongsToHousehold = input.belongsToHousehold = '';
    this.createdAt = input.createdAt = '';
  }
}
