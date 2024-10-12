// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidMeasurementUnitConversionUpdateRequestInput {
   from: string;
 modifier: number;
 notes: string;
 onlyForIngredient: string;
 to: string;

}

export class ValidMeasurementUnitConversionUpdateRequestInput implements IValidMeasurementUnitConversionUpdateRequestInput {
   from: string;
 modifier: number;
 notes: string;
 onlyForIngredient: string;
 to: string;
constructor(input: Partial<ValidMeasurementUnitConversionUpdateRequestInput> = {}) {
	 this.from = input.from || '';
 this.modifier = input.modifier || 0;
 this.notes = input.notes || '';
 this.onlyForIngredient = input.onlyForIngredient || '';
 this.to = input.to || '';
}
}