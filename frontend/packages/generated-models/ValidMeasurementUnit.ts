// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IValidMeasurementUnit {
   name: string;
 slug: string;
 volumetric: boolean;
 archivedAt?: string;
 createdAt: string;
 iconPath: string;
 metric: boolean;
 pluralName: string;
 universal: boolean;
 description: string;
 id: string;
 imperial: boolean;
 lastUpdatedAt?: string;

}

export class ValidMeasurementUnit implements IValidMeasurementUnit {
   name: string;
 slug: string;
 volumetric: boolean;
 archivedAt?: string;
 createdAt: string;
 iconPath: string;
 metric: boolean;
 pluralName: string;
 universal: boolean;
 description: string;
 id: string;
 imperial: boolean;
 lastUpdatedAt?: string;
constructor(input: Partial<ValidMeasurementUnit> = {}) {
	 this.name = input.name = '';
 this.slug = input.slug = '';
 this.volumetric = input.volumetric = false;
 this.archivedAt = input.archivedAt;
 this.createdAt = input.createdAt = '';
 this.iconPath = input.iconPath = '';
 this.metric = input.metric = false;
 this.pluralName = input.pluralName = '';
 this.universal = input.universal = false;
 this.description = input.description = '';
 this.id = input.id = '';
 this.imperial = input.imperial = false;
 this.lastUpdatedAt = input.lastUpdatedAt;
}
}