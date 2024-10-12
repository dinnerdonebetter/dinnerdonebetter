// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidVesselShapeType } from './enums';


export interface IValidVesselUpdateRequestInput {
   capacity: number;
 capacityUnitID: string;
 description: string;
 displayInSummaryLists: boolean;
 heightInMillimeters: number;
 iconPath: string;
 includeInGeneratedInstructions: boolean;
 lengthInMillimeters: number;
 name: string;
 pluralName: string;
 shape: ValidVesselShapeType;
 slug: string;
 usableForStorage: boolean;
 widthInMillimeters: number;

}

export class ValidVesselUpdateRequestInput implements IValidVesselUpdateRequestInput {
   capacity: number;
 capacityUnitID: string;
 description: string;
 displayInSummaryLists: boolean;
 heightInMillimeters: number;
 iconPath: string;
 includeInGeneratedInstructions: boolean;
 lengthInMillimeters: number;
 name: string;
 pluralName: string;
 shape: ValidVesselShapeType;
 slug: string;
 usableForStorage: boolean;
 widthInMillimeters: number;
constructor(input: Partial<ValidVesselUpdateRequestInput> = {}) {
	 this.capacity = input.capacity || 0;
 this.capacityUnitID = input.capacityUnitID || '';
 this.description = input.description || '';
 this.displayInSummaryLists = input.displayInSummaryLists || false;
 this.heightInMillimeters = input.heightInMillimeters || 0;
 this.iconPath = input.iconPath || '';
 this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
 this.lengthInMillimeters = input.lengthInMillimeters || 0;
 this.name = input.name || '';
 this.pluralName = input.pluralName || '';
 this.shape = input.shape || 'other';
 this.slug = input.slug || '';
 this.usableForStorage = input.usableForStorage || false;
 this.widthInMillimeters = input.widthInMillimeters || 0;
}
}