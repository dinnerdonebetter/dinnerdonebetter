// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidIngredient } from './ValidIngredient';
 import { ValidMeasurementUnit } from './ValidMeasurementUnit';


export interface IValidMeasurementUnitConversion {
   from: ValidMeasurementUnit;
 modifier: number;
 notes: string;
 onlyForIngredient?: ValidIngredient;
 archivedAt?: string;
 createdAt: string;
 id: string;
 lastUpdatedAt?: string;
 to: ValidMeasurementUnit;

}

export class ValidMeasurementUnitConversion implements IValidMeasurementUnitConversion {
   from: ValidMeasurementUnit;
 modifier: number;
 notes: string;
 onlyForIngredient?: ValidIngredient;
 archivedAt?: string;
 createdAt: string;
 id: string;
 lastUpdatedAt?: string;
 to: ValidMeasurementUnit;
constructor(input: Partial<ValidMeasurementUnitConversion> = {}) {
	 this.from = input.from = new ValidMeasurementUnit();
 this.modifier = input.modifier = 0;
 this.notes = input.notes = '';
 this.onlyForIngredient = input.onlyForIngredient;
 this.archivedAt = input.archivedAt;
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.to = input.to = new ValidMeasurementUnit();
}
}