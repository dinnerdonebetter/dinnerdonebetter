// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidMeasurementUnitConversionUpdateRequestInput {
   notes?: string;
 onlyForIngredient?: string;
 to?: string;
 from?: string;
 modifier?: number;

}

export class ValidMeasurementUnitConversionUpdateRequestInput implements IValidMeasurementUnitConversionUpdateRequestInput {
   notes?: string;
 onlyForIngredient?: string;
 to?: string;
 from?: string;
 modifier?: number;
constructor(input: Partial<ValidMeasurementUnitConversionUpdateRequestInput> = {}) {
	 this.notes = input.notes;
 this.onlyForIngredient = input.onlyForIngredient;
 this.to = input.to;
 this.from = input.from;
 this.modifier = input.modifier;
}
}