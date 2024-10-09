// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidMeasurementUnitCreationRequestInput {
   iconPath: string;
 imperial: boolean;
 metric: boolean;
 name: string;
 pluralName: string;
 slug: string;
 volumetric: boolean;
 description: string;
 universal: boolean;

}

export class ValidMeasurementUnitCreationRequestInput implements IValidMeasurementUnitCreationRequestInput {
   iconPath: string;
 imperial: boolean;
 metric: boolean;
 name: string;
 pluralName: string;
 slug: string;
 volumetric: boolean;
 description: string;
 universal: boolean;
constructor(input: Partial<ValidMeasurementUnitCreationRequestInput> = {}) {
	 this.iconPath = input.iconPath = '';
 this.imperial = input.imperial = false;
 this.metric = input.metric = false;
 this.name = input.name = '';
 this.pluralName = input.pluralName = '';
 this.slug = input.slug = '';
 this.volumetric = input.volumetric = false;
 this.description = input.description = '';
 this.universal = input.universal = false;
}
}