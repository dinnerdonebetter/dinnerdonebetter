// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidInstrument } from './ValidInstrument.gen';
import { ValidPreparation } from './ValidPreparation.gen';

export interface IValidPreparationInstrument {
  archivedAt: string;
  createdAt: string;
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
}

export class ValidPreparationInstrument implements IValidPreparationInstrument {
  archivedAt: string;
  createdAt: string;
  id: string;
  instrument: ValidInstrument;
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
  constructor(input: Partial<ValidPreparationInstrument> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.instrument = input.instrument || new ValidInstrument();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.preparation = input.preparation || new ValidPreparation();
  }
}
