// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidMeasurementUnitCreationRequestInput {
   description: string;
 iconPath: string;
 imperial: boolean;
 metric: boolean;
 name: string;
 pluralName: string;
 slug: string;
 universal: boolean;
 volumetric: boolean;

}

export class ValidMeasurementUnitCreationRequestInput implements IValidMeasurementUnitCreationRequestInput {
   description: string;
 iconPath: string;
 imperial: boolean;
 metric: boolean;
 name: string;
 pluralName: string;
 slug: string;
 universal: boolean;
 volumetric: boolean;
constructor(input: Partial<ValidMeasurementUnitCreationRequestInput> = {}) {
	 this.description = input.description || '';
 this.iconPath = input.iconPath || '';
 this.imperial = input.imperial || false;
 this.metric = input.metric || false;
 this.name = input.name || '';
 this.pluralName = input.pluralName || '';
 this.slug = input.slug || '';
 this.universal = input.universal || false;
 this.volumetric = input.volumetric || false;
}
}