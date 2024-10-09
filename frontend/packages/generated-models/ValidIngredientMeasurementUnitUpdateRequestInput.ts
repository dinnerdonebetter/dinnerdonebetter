// GENERATED CODE, DO NOT EDIT MANUALLY

import { OptionalNumberRange } from './number_range';

export interface IValidIngredientMeasurementUnitUpdateRequestInput {
  notes?: string;
  validIngredientID?: string;
  validMeasurementUnitID?: string;
  allowableQuantity: OptionalNumberRange;
}

export class ValidIngredientMeasurementUnitUpdateRequestInput
  implements IValidIngredientMeasurementUnitUpdateRequestInput
{
  notes?: string;
  validIngredientID?: string;
  validMeasurementUnitID?: string;
  allowableQuantity: OptionalNumberRange;
  constructor(input: Partial<ValidIngredientMeasurementUnitUpdateRequestInput> = {}) {
    this.notes = input.notes;
    this.validIngredientID = input.validIngredientID;
    this.validMeasurementUnitID = input.validMeasurementUnitID;
    this.allowableQuantity = input.allowableQuantity = {};
  }
}
