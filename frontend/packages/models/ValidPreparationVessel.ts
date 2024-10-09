// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation';
import { ValidVessel } from './ValidVessel';

export interface IValidPreparationVessel {
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
  archivedAt: string;
}

export class ValidPreparationVessel implements IValidPreparationVessel {
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
  archivedAt: string;
  constructor(input: Partial<ValidPreparationVessel> = {}) {
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.instrument = input.instrument || new ValidVessel();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.preparation = input.preparation || new ValidPreparation();
    this.archivedAt = input.archivedAt || '';
  }
}
