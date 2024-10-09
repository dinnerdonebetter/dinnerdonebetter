// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVesselShapeType } from './enums';

export interface IValidVesselCreationRequestInput {
  displayInSummaryLists: boolean;
  description: string;
  usableForStorage: boolean;
  capacityUnitID: string;
  iconPath: string;
  shape: ValidVesselShapeType;
  slug: string;
  capacity: number;
  heightInMillimeters: number;
  includeInGeneratedInstructions: boolean;
  lengthInMillimeters: number;
  name: string;
  pluralName: string;
  widthInMillimeters: number;
}

export class ValidVesselCreationRequestInput implements IValidVesselCreationRequestInput {
  displayInSummaryLists: boolean;
  description: string;
  usableForStorage: boolean;
  capacityUnitID: string;
  iconPath: string;
  shape: ValidVesselShapeType;
  slug: string;
  capacity: number;
  heightInMillimeters: number;
  includeInGeneratedInstructions: boolean;
  lengthInMillimeters: number;
  name: string;
  pluralName: string;
  widthInMillimeters: number;
  constructor(input: Partial<ValidVesselCreationRequestInput> = {}) {
    this.displayInSummaryLists = input.displayInSummaryLists || false;
    this.description = input.description || '';
    this.usableForStorage = input.usableForStorage || false;
    this.capacityUnitID = input.capacityUnitID || '';
    this.iconPath = input.iconPath || '';
    this.shape = input.shape || 'other';
    this.slug = input.slug || '';
    this.capacity = input.capacity || 0;
    this.heightInMillimeters = input.heightInMillimeters || 0;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
    this.lengthInMillimeters = input.lengthInMillimeters || 0;
    this.name = input.name || '';
    this.pluralName = input.pluralName || '';
    this.widthInMillimeters = input.widthInMillimeters || 0;
  }
}
