// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitConversionCreationRequestInput {
  from: string;
  modifier: number;
  notes: string;
  onlyForIngredient: string;
  to: string;
}

export class ValidMeasurementUnitConversionCreationRequestInput
  implements IValidMeasurementUnitConversionCreationRequestInput
{
  from: string;
  modifier: number;
  notes: string;
  onlyForIngredient: string;
  to: string;
  constructor(input: Partial<ValidMeasurementUnitConversionCreationRequestInput> = {}) {
    this.from = input.from || '';
    this.modifier = input.modifier || 0;
    this.notes = input.notes || '';
    this.onlyForIngredient = input.onlyForIngredient || '';
    this.to = input.to || '';
  }
}
