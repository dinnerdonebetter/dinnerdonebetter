// GENERATED CODE, DO NOT EDIT MANUALLY

 import { NumberRangeWithOptionalMax } from './number_range';


export interface IValidIngredientMeasurementUnitCreationRequestInput {
   allowableQuantity: NumberRangeWithOptionalMax;
 notes: string;
 validIngredientID: string;
 validMeasurementUnitID: string;

}

export class ValidIngredientMeasurementUnitCreationRequestInput implements IValidIngredientMeasurementUnitCreationRequestInput {
   allowableQuantity: NumberRangeWithOptionalMax;
 notes: string;
 validIngredientID: string;
 validMeasurementUnitID: string;
constructor(input: Partial<ValidIngredientMeasurementUnitCreationRequestInput> = {}) {
	 this.allowableQuantity = input.allowableQuantity || { min: 0 };
 this.notes = input.notes || '';
 this.validIngredientID = input.validIngredientID || '';
 this.validMeasurementUnitID = input.validMeasurementUnitID || '';
}
}