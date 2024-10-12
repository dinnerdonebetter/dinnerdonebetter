// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidMeasurementUnit } from './ValidMeasurementUnit';
 import { ValidVesselShapeType } from './enums';


export interface IValidVessel {
   archivedAt: string;
 capacity: number;
 capacityUnit: ValidMeasurementUnit;
 createdAt: string;
 description: string;
 displayInSummaryLists: boolean;
 heightInMillimeters: number;
 iconPath: string;
 id: string;
 includeInGeneratedInstructions: boolean;
 lastUpdatedAt: string;
 lengthInMillimeters: number;
 name: string;
 pluralName: string;
 shape: ValidVesselShapeType;
 slug: string;
 usableForStorage: boolean;
 widthInMillimeters: number;

}

export class ValidVessel implements IValidVessel {
   archivedAt: string;
 capacity: number;
 capacityUnit: ValidMeasurementUnit;
 createdAt: string;
 description: string;
 displayInSummaryLists: boolean;
 heightInMillimeters: number;
 iconPath: string;
 id: string;
 includeInGeneratedInstructions: boolean;
 lastUpdatedAt: string;
 lengthInMillimeters: number;
 name: string;
 pluralName: string;
 shape: ValidVesselShapeType;
 slug: string;
 usableForStorage: boolean;
 widthInMillimeters: number;
constructor(input: Partial<ValidVessel> = {}) {
	 this.archivedAt = input.archivedAt || '';
 this.capacity = input.capacity || 0;
 this.capacityUnit = input.capacityUnit || new ValidMeasurementUnit();
 this.createdAt = input.createdAt || '';
 this.description = input.description || '';
 this.displayInSummaryLists = input.displayInSummaryLists || false;
 this.heightInMillimeters = input.heightInMillimeters || 0;
 this.iconPath = input.iconPath || '';
 this.id = input.id || '';
 this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.lengthInMillimeters = input.lengthInMillimeters || 0;
 this.name = input.name || '';
 this.pluralName = input.pluralName || '';
 this.shape = input.shape || 'other';
 this.slug = input.slug || '';
 this.usableForStorage = input.usableForStorage || false;
 this.widthInMillimeters = input.widthInMillimeters || 0;
}
}