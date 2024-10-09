// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitConversionCreationRequestInput {
  onlyForIngredient?: string;
  to: string;
  from: string;
  modifier: number;
  notes: string;
}

export class ValidMeasurementUnitConversionCreationRequestInput
  implements IValidMeasurementUnitConversionCreationRequestInput
{
  onlyForIngredient?: string;
  to: string;
  from: string;
  modifier: number;
  notes: string;
  constructor(input: Partial<ValidMeasurementUnitConversionCreationRequestInput> = {}) {
    this.onlyForIngredient = input.onlyForIngredient;
    this.to = input.to = '';
    this.from = input.from = '';
    this.modifier = input.modifier = 0;
    this.notes = input.notes = '';
  }
}
