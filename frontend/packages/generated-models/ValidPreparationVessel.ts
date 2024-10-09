// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation';
import { ValidVessel } from './ValidVessel';

export interface IValidPreparationVessel {
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  lastUpdatedAt?: string;
  notes: string;
}

export class ValidPreparationVessel implements IValidPreparationVessel {
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  lastUpdatedAt?: string;
  notes: string;
  constructor(input: Partial<ValidPreparationVessel> = {}) {
    this.preparation = input.preparation = new ValidPreparation();
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.instrument = input.instrument = new ValidVessel();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
  }
}
