// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVesselShapeType } from './enums';

export interface IValidVesselUpdateRequestInput {
  description: string;
  displayInSummaryLists: boolean;
  includeInGeneratedInstructions: boolean;
  widthInMillimeters: number;
  capacityUnitID: string;
  heightInMillimeters: number;
  slug: string;
  capacity: number;
  iconPath: string;
  name: string;
  pluralName: string;
  usableForStorage: boolean;
  lengthInMillimeters: number;
  shape: ValidVesselShapeType;
}

export class ValidVesselUpdateRequestInput implements IValidVesselUpdateRequestInput {
  description: string;
  displayInSummaryLists: boolean;
  includeInGeneratedInstructions: boolean;
  widthInMillimeters: number;
  capacityUnitID: string;
  heightInMillimeters: number;
  slug: string;
  capacity: number;
  iconPath: string;
  name: string;
  pluralName: string;
  usableForStorage: boolean;
  lengthInMillimeters: number;
  shape: ValidVesselShapeType;
  constructor(input: Partial<ValidVesselUpdateRequestInput> = {}) {
    this.description = input.description || '';
    this.displayInSummaryLists = input.displayInSummaryLists || false;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
    this.widthInMillimeters = input.widthInMillimeters || 0;
    this.capacityUnitID = input.capacityUnitID || '';
    this.heightInMillimeters = input.heightInMillimeters || 0;
    this.slug = input.slug || '';
    this.capacity = input.capacity || 0;
    this.iconPath = input.iconPath || '';
    this.name = input.name || '';
    this.pluralName = input.pluralName || '';
    this.usableForStorage = input.usableForStorage || false;
    this.lengthInMillimeters = input.lengthInMillimeters || 0;
    this.shape = input.shape || 'other';
  }
}
