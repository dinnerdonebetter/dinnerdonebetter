// GENERATED CODE, DO NOT EDIT MANUALLY

 import { ValidVesselShapeType } from './enums';


export interface IValidVesselUpdateRequestInput {
   slug?: string;
 capacityUnitID?: string;
 includeInGeneratedInstructions?: boolean;
 lengthInMillimeters?: number;
 shape?: ValidVesselShapeType;
 capacity?: number;
 displayInSummaryLists?: boolean;
 iconPath?: string;
 usableForStorage?: boolean;
 description?: string;
 widthInMillimeters?: number;
 heightInMillimeters?: number;
 name?: string;
 pluralName?: string;

}

export class ValidVesselUpdateRequestInput implements IValidVesselUpdateRequestInput {
   slug?: string;
 capacityUnitID?: string;
 includeInGeneratedInstructions?: boolean;
 lengthInMillimeters?: number;
 shape?: ValidVesselShapeType;
 capacity?: number;
 displayInSummaryLists?: boolean;
 iconPath?: string;
 usableForStorage?: boolean;
 description?: string;
 widthInMillimeters?: number;
 heightInMillimeters?: number;
 name?: string;
 pluralName?: string;
constructor(input: Partial<ValidVesselUpdateRequestInput> = {}) {
	 this.slug = input.slug;
 this.capacityUnitID = input.capacityUnitID;
 this.includeInGeneratedInstructions = input.includeInGeneratedInstructions;
 this.lengthInMillimeters = input.lengthInMillimeters;
 this.shape = input.shape;
 this.capacity = input.capacity;
 this.displayInSummaryLists = input.displayInSummaryLists;
 this.iconPath = input.iconPath;
 this.usableForStorage = input.usableForStorage;
 this.description = input.description;
 this.widthInMillimeters = input.widthInMillimeters;
 this.heightInMillimeters = input.heightInMillimeters;
 this.name = input.name;
 this.pluralName = input.pluralName;
}
}