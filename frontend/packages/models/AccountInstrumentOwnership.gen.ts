// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument.gen';

export interface IAccountInstrumentOwnership {
  archivedAt: string;
  belongsToAccount: string;
  createdAt: string;
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt: string;
  notes: string;
  quantity: number;
}

export class AccountInstrumentOwnership implements IAccountInstrumentOwnership {
  archivedAt: string;
  belongsToAccount: string;
  createdAt: string;
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt: string;
  notes: string;
  quantity: number;
  constructor(input: Partial<AccountInstrumentOwnership> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToAccount = input.belongsToAccount || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.instrument = input.instrument || new ValidInstrument();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.quantity = input.quantity || 0;
  }
}
