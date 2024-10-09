// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IValidIngredientMeasurementUnitUpdateRequestInput {
  validMeasurementUnitID: string;
  allowableQuantity: OptionalNumberRange;
  notes: string;
  validIngredientID: string;
}

export class ValidIngredientMeasurementUnitUpdateRequestInput
  implements IValidIngredientMeasurementUnitUpdateRequestInput
{
  validMeasurementUnitID: string;
  allowableQuantity: OptionalNumberRange;
  notes: string;
  validIngredientID: string;
  constructor(input: Partial<ValidIngredientMeasurementUnitUpdateRequestInput> = {}) {
    this.validMeasurementUnitID = input.validMeasurementUnitID || '';
    this.allowableQuantity = input.allowableQuantity || {};
    this.notes = input.notes || '';
    this.validIngredientID = input.validIngredientID || '';
  }
}
