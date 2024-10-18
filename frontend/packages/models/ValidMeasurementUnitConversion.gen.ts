// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidIngredient } from './ValidIngredient.gen';
import { ValidMeasurementUnit } from './ValidMeasurementUnit.gen';

export interface IValidMeasurementUnitConversion {
  archivedAt: string;
  createdAt: string;
  from: ValidMeasurementUnit;
  id: string;
  lastUpdatedAt: string;
  modifier: number;
  notes: string;
  onlyForIngredient: ValidIngredient;
  to: ValidMeasurementUnit;
}

export class ValidMeasurementUnitConversion implements IValidMeasurementUnitConversion {
  archivedAt: string;
  createdAt: string;
  from: ValidMeasurementUnit;
  id: string;
  lastUpdatedAt: string;
  modifier: number;
  notes: string;
  onlyForIngredient: ValidIngredient;
  to: ValidMeasurementUnit;
  constructor(input: Partial<ValidMeasurementUnitConversion> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.from = input.from || new ValidMeasurementUnit();
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.modifier = input.modifier || 0;
    this.notes = input.notes || '';
    this.onlyForIngredient = input.onlyForIngredient || new ValidIngredient();
    this.to = input.to || new ValidMeasurementUnit();
  }
}
