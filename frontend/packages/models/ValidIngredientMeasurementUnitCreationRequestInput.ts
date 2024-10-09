// GENERATED CODE, DO NOT EDIT MANUALLY

import { NumberRangeWithOptionalMax } from './number_range';

export interface IValidIngredientMeasurementUnitCreationRequestInput {
  validMeasurementUnitID: string;
  allowableQuantity: NumberRangeWithOptionalMax;
  notes: string;
  validIngredientID: string;
}

export class ValidIngredientMeasurementUnitCreationRequestInput
  implements IValidIngredientMeasurementUnitCreationRequestInput
{
  validMeasurementUnitID: string;
  allowableQuantity: NumberRangeWithOptionalMax;
  notes: string;
  validIngredientID: string;
  constructor(input: Partial<ValidIngredientMeasurementUnitCreationRequestInput> = {}) {
    this.validMeasurementUnitID = input.validMeasurementUnitID || '';
    this.allowableQuantity = input.allowableQuantity || { min: 0 };
    this.notes = input.notes || '';
    this.validIngredientID = input.validIngredientID || '';
  }
}
