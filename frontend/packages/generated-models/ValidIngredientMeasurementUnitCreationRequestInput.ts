// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidIngredientMeasurementUnitCreationRequestInput {
  notes: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
  allowableQuantity: NumberRangeWithOptionalMax;
}

export class ValidIngredientMeasurementUnitCreationRequestInput
  implements IValidIngredientMeasurementUnitCreationRequestInput
{
  notes: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
  allowableQuantity: NumberRangeWithOptionalMax;
  constructor(input: Partial<ValidIngredientMeasurementUnitCreationRequestInput> = {}) {
    this.notes = input.notes = '';
    this.validIngredientID = input.validIngredientID = '';
    this.validMeasurementUnitID = input.validMeasurementUnitID = '';
    this.allowableQuantity = input.allowableQuantity = { min: 0 };
  }
}
