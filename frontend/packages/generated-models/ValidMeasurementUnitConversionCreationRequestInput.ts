// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitConversionCreationRequestInput {
  notes: string;
  onlyForIngredient?: string;
  to: string;
  from: string;
  modifier: number;
}

export class ValidMeasurementUnitConversionCreationRequestInput
  implements IValidMeasurementUnitConversionCreationRequestInput
{
  notes: string;
  onlyForIngredient?: string;
  to: string;
  from: string;
  modifier: number;
  constructor(input: Partial<ValidMeasurementUnitConversionCreationRequestInput> = {}) {
    this.notes = input.notes = '';
    this.onlyForIngredient = input.onlyForIngredient;
    this.to = input.to = '';
    this.from = input.from = '';
    this.modifier = input.modifier = 0;
  }
}
