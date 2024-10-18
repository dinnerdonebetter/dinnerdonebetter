// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidPreparation } from './ValidPreparation.gen';
import { ValidVessel } from './ValidVessel.gen';

export interface IValidPreparationVessel {
  archivedAt: string;
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
}

export class ValidPreparationVessel implements IValidPreparationVessel {
  archivedAt: string;
  createdAt: string;
  id: string;
  instrument: ValidVessel;
  lastUpdatedAt: string;
  notes: string;
  preparation: ValidPreparation;
  constructor(input: Partial<ValidPreparationVessel> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.instrument = input.instrument || new ValidVessel();
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.notes = input.notes || '';
    this.preparation = input.preparation || new ValidPreparation();
  }
}
