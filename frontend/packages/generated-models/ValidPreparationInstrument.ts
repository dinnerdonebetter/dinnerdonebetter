// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument';
import { ValidPreparation } from './ValidPreparation';

export interface IValidPreparationInstrument {
  notes: string;
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt?: string;
}

export class ValidPreparationInstrument implements IValidPreparationInstrument {
  notes: string;
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt?: string;
  constructor(input: Partial<ValidPreparationInstrument> = {}) {
    this.notes = input.notes = '';
    this.preparation = input.preparation = new ValidPreparation();
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.instrument = input.instrument = new ValidInstrument();
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
