// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidVesselShapeType } from './enums';


export interface IValidVesselCreationRequestInput {
   heightInMillimeters: number;
 capacityUnitID?: string;
 usableForStorage: boolean;
 widthInMillimeters: number;
 slug: string;
 iconPath: string;
 lengthInMillimeters: number;
 pluralName: string;
 description: string;
 displayInSummaryLists: boolean;
 includeInGeneratedInstructions: boolean;
 name: string;
 shape: ValidVesselShapeType;
 capacity: number;

}

export class ValidVesselCreationRequestInput implements IValidVesselCreationRequestInput {
   heightInMillimeters: number;
 capacityUnitID?: string;
 usableForStorage: boolean;
 widthInMillimeters: number;
 slug: string;
 iconPath: string;
 lengthInMillimeters: number;
 pluralName: string;
 description: string;
 displayInSummaryLists: boolean;
 includeInGeneratedInstructions: boolean;
 name: string;
 shape: ValidVesselShapeType;
 capacity: number;
constructor(input: Partial<ValidVesselCreationRequestInput> = {}) {
	 this.heightInMillimeters = input.heightInMillimeters = 0;
 this.capacityUnitID = input.capacityUnitID;
 this.usableForStorage = input.usableForStorage = false;
 this.widthInMillimeters = input.widthInMillimeters = 0;
 this.slug = input.slug = '';
 this.iconPath = input.iconPath = '';
 this.lengthInMillimeters = input.lengthInMillimeters = 0;
 this.pluralName = input.pluralName = '';
 this.description = input.description = '';
 this.displayInSummaryLists = input.displayInSummaryLists = false;
 this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
 this.name = input.name = '';
 this.shape = input.shape = 'other';
 this.capacity = input.capacity = 0;
}
}