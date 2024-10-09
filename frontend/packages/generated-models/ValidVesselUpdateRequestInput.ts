// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVesselShapeType } from './enums';

export interface IValidVesselUpdateRequestInput {
  capacity?: number;
  includeInGeneratedInstructions?: boolean;
  widthInMillimeters?: number;
  displayInSummaryLists?: boolean;
  iconPath?: string;
  lengthInMillimeters?: number;
  name?: string;
  pluralName?: string;
  capacityUnitID?: string;
  heightInMillimeters?: number;
  shape?: ValidVesselShapeType;
  usableForStorage?: boolean;
  description?: string;
  slug?: string;
}

export class ValidVesselUpdateRequestInput implements IValidVesselUpdateRequestInput {
  capacity?: number;
  includeInGeneratedInstructions?: boolean;
  widthInMillimeters?: number;
  displayInSummaryLists?: boolean;
  iconPath?: string;
  lengthInMillimeters?: number;
  name?: string;
  pluralName?: string;
  capacityUnitID?: string;
  heightInMillimeters?: number;
  shape?: ValidVesselShapeType;
  usableForStorage?: boolean;
  description?: string;
  slug?: string;
  constructor(input: Partial<ValidVesselUpdateRequestInput> = {}) {
    this.capacity = input.capacity;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions;
    this.widthInMillimeters = input.widthInMillimeters;
    this.displayInSummaryLists = input.displayInSummaryLists;
    this.iconPath = input.iconPath;
    this.lengthInMillimeters = input.lengthInMillimeters;
    this.name = input.name;
    this.pluralName = input.pluralName;
    this.capacityUnitID = input.capacityUnitID;
    this.heightInMillimeters = input.heightInMillimeters;
    this.shape = input.shape;
    this.usableForStorage = input.usableForStorage;
    this.description = input.description;
    this.slug = input.slug;
  }
}
