// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitConversionCreationRequestInput {
  to: string;
  from: string;
  modifier: number;
  notes: string;
  onlyForIngredient: string;
}

export class ValidMeasurementUnitConversionCreationRequestInput
  implements IValidMeasurementUnitConversionCreationRequestInput
{
  to: string;
  from: string;
  modifier: number;
  notes: string;
  onlyForIngredient: string;
  constructor(input: Partial<ValidMeasurementUnitConversionCreationRequestInput> = {}) {
    this.to = input.to || '';
    this.from = input.from || '';
    this.modifier = input.modifier || 0;
    this.notes = input.notes || '';
    this.onlyForIngredient = input.onlyForIngredient || '';
  }
}
