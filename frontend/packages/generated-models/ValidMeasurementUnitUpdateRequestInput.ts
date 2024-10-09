// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidMeasurementUnitUpdateRequestInput {
   imperial?: boolean;
 pluralName?: string;
 slug?: string;
 volumetric?: boolean;
 description?: string;
 metric?: boolean;
 name?: string;
 universal?: boolean;
 iconPath?: string;

}

export class ValidMeasurementUnitUpdateRequestInput implements IValidMeasurementUnitUpdateRequestInput {
   imperial?: boolean;
 pluralName?: string;
 slug?: string;
 volumetric?: boolean;
 description?: string;
 metric?: boolean;
 name?: string;
 universal?: boolean;
 iconPath?: string;
constructor(input: Partial<ValidMeasurementUnitUpdateRequestInput> = {}) {
	 this.imperial = input.imperial;
 this.pluralName = input.pluralName;
 this.slug = input.slug;
 this.volumetric = input.volumetric;
 this.description = input.description;
 this.metric = input.metric;
 this.name = input.name;
 this.universal = input.universal;
 this.iconPath = input.iconPath;
}
}