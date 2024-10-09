// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation';
import { ValidVessel } from './ValidVessel';

export interface IValidPreparationVessel {
  notes: string;
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  lastUpdatedAt?: string;
}

export class ValidPreparationVessel implements IValidPreparationVessel {
  notes: string;
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  lastUpdatedAt?: string;
  constructor(input: Partial<ValidPreparationVessel> = {}) {
    this.notes = input.notes = '';
    this.preparation = input.preparation = new ValidPreparation();
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.instrument = input.instrument = new ValidVessel();
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
