// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidMeasurementUnit } from './ValidMeasurementUnit';
 import { ValidVesselShapeType } from './enums';


export interface IValidVessel {
   iconPath: string;
 capacityUnit?: ValidMeasurementUnit;
 createdAt: string;
 heightInMillimeters: number;
 id: string;
 includeInGeneratedInstructions: boolean;
 shape: ValidVesselShapeType;
 archivedAt?: string;
 displayInSummaryLists: boolean;
 name: string;
 slug: string;
 usableForStorage: boolean;
 capacity: number;
 description: string;
 lastUpdatedAt?: string;
 lengthInMillimeters: number;
 pluralName: string;
 widthInMillimeters: number;

}

export class ValidVessel implements IValidVessel {
   iconPath: string;
 capacityUnit?: ValidMeasurementUnit;
 createdAt: string;
 heightInMillimeters: number;
 id: string;
 includeInGeneratedInstructions: boolean;
 shape: ValidVesselShapeType;
 archivedAt?: string;
 displayInSummaryLists: boolean;
 name: string;
 slug: string;
 usableForStorage: boolean;
 capacity: number;
 description: string;
 lastUpdatedAt?: string;
 lengthInMillimeters: number;
 pluralName: string;
 widthInMillimeters: number;
constructor(input: Partial<ValidVessel> = {}) {
	 this.iconPath = input.iconPath = '';
 this.capacityUnit = input.capacityUnit;
 this.createdAt = input.createdAt = '';
 this.heightInMillimeters = input.heightInMillimeters = 0;
 this.id = input.id = '';
 this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
 this.shape = input.shape = 'other';
 this.archivedAt = input.archivedAt;
 this.displayInSummaryLists = input.displayInSummaryLists = false;
 this.name = input.name = '';
 this.slug = input.slug = '';
 this.usableForStorage = input.usableForStorage = false;
 this.capacity = input.capacity = 0;
 this.description = input.description = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.lengthInMillimeters = input.lengthInMillimeters = 0;
 this.pluralName = input.pluralName = '';
 this.widthInMillimeters = input.widthInMillimeters = 0;
}
}