// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation';
import { ValidVessel } from './ValidVessel';

export interface IValidPreparationVessel {
  lastUpdatedAt?: string;
  notes: string;
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  id: string;
  instrument: ValidVessel;
}

export class ValidPreparationVessel implements IValidPreparationVessel {
  lastUpdatedAt?: string;
  notes: string;
  preparation: ValidPreparation;
  archivedAt?: string;
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  constructor(input: Partial<ValidPreparationVessel> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.notes = input.notes = '';
    this.preparation = input.preparation = new ValidPreparation();
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.instrument = input.instrument = new ValidVessel();
  }
}
