// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitConversionUpdateRequestInput {
  to: string;
  from: string;
  modifier: number;
  notes: string;
  onlyForIngredient: string;
}

export class ValidMeasurementUnitConversionUpdateRequestInput
  implements IValidMeasurementUnitConversionUpdateRequestInput
{
  to: string;
  from: string;
  modifier: number;
  notes: string;
  onlyForIngredient: string;
  constructor(input: Partial<ValidMeasurementUnitConversionUpdateRequestInput> = {}) {
    this.to = input.to || '';
    this.from = input.from || '';
    this.modifier = input.modifier || 0;
    this.notes = input.notes || '';
    this.onlyForIngredient = input.onlyForIngredient || '';
  }
}
