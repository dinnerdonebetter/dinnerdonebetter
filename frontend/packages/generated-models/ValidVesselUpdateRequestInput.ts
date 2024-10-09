// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVesselShapeType } from './enums';

export interface IValidVesselUpdateRequestInput {
  lengthInMillimeters?: number;
  usableForStorage?: boolean;
  name?: string;
  pluralName?: string;
  shape?: ValidVesselShapeType;
  capacity?: number;
  displayInSummaryLists?: boolean;
  heightInMillimeters?: number;
  iconPath?: string;
  widthInMillimeters?: number;
  capacityUnitID?: string;
  description?: string;
  includeInGeneratedInstructions?: boolean;
  slug?: string;
}

export class ValidVesselUpdateRequestInput implements IValidVesselUpdateRequestInput {
  lengthInMillimeters?: number;
  usableForStorage?: boolean;
  name?: string;
  pluralName?: string;
  shape?: ValidVesselShapeType;
  capacity?: number;
  displayInSummaryLists?: boolean;
  heightInMillimeters?: number;
  iconPath?: string;
  widthInMillimeters?: number;
  capacityUnitID?: string;
  description?: string;
  includeInGeneratedInstructions?: boolean;
  slug?: string;
  constructor(input: Partial<ValidVesselUpdateRequestInput> = {}) {
    this.lengthInMillimeters = input.lengthInMillimeters;
    this.usableForStorage = input.usableForStorage;
    this.name = input.name;
    this.pluralName = input.pluralName;
    this.shape = input.shape;
    this.capacity = input.capacity;
    this.displayInSummaryLists = input.displayInSummaryLists;
    this.heightInMillimeters = input.heightInMillimeters;
    this.iconPath = input.iconPath;
    this.widthInMillimeters = input.widthInMillimeters;
    this.capacityUnitID = input.capacityUnitID;
    this.description = input.description;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions;
    this.slug = input.slug;
  }
}
