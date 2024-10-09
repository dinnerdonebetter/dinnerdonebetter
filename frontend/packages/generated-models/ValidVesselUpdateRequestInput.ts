// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVesselShapeType } from './enums';

export interface IValidVesselUpdateRequestInput {
  displayInSummaryLists?: boolean;
  heightInMillimeters?: number;
  description?: string;
  capacity?: number;
  iconPath?: string;
  includeInGeneratedInstructions?: boolean;
  name?: string;
  shape?: ValidVesselShapeType;
  usableForStorage?: boolean;
  capacityUnitID?: string;
  lengthInMillimeters?: number;
  pluralName?: string;
  slug?: string;
  widthInMillimeters?: number;
}

export class ValidVesselUpdateRequestInput implements IValidVesselUpdateRequestInput {
  displayInSummaryLists?: boolean;
  heightInMillimeters?: number;
  description?: string;
  capacity?: number;
  iconPath?: string;
  includeInGeneratedInstructions?: boolean;
  name?: string;
  shape?: ValidVesselShapeType;
  usableForStorage?: boolean;
  capacityUnitID?: string;
  lengthInMillimeters?: number;
  pluralName?: string;
  slug?: string;
  widthInMillimeters?: number;
  constructor(input: Partial<ValidVesselUpdateRequestInput> = {}) {
    this.displayInSummaryLists = input.displayInSummaryLists;
    this.heightInMillimeters = input.heightInMillimeters;
    this.description = input.description;
    this.capacity = input.capacity;
    this.iconPath = input.iconPath;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions;
    this.name = input.name;
    this.shape = input.shape;
    this.usableForStorage = input.usableForStorage;
    this.capacityUnitID = input.capacityUnitID;
    this.lengthInMillimeters = input.lengthInMillimeters;
    this.pluralName = input.pluralName;
    this.slug = input.slug;
    this.widthInMillimeters = input.widthInMillimeters;
  }
}
