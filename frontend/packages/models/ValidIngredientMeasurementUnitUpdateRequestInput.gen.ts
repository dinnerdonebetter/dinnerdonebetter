// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IValidIngredientMeasurementUnitUpdateRequestInput {
  allowableQuantity: OptionalNumberRange;
  notes: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
}

export class ValidIngredientMeasurementUnitUpdateRequestInput
  implements IValidIngredientMeasurementUnitUpdateRequestInput
{
  allowableQuantity: OptionalNumberRange;
  notes: string;
  validIngredientID: string;
  validMeasurementUnitID: string;
  constructor(input: Partial<ValidIngredientMeasurementUnitUpdateRequestInput> = {}) {
    this.allowableQuantity = input.allowableQuantity || {};
    this.notes = input.notes || '';
    this.validIngredientID = input.validIngredientID || '';
    this.validMeasurementUnitID = input.validMeasurementUnitID || '';
  }
}
