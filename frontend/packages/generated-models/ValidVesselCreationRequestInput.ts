// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVesselShapeType } from './enums';

export interface IValidVesselCreationRequestInput {
  capacityUnitID?: string;
  displayInSummaryLists: boolean;
  lengthInMillimeters: number;
  description: string;
  iconPath: string;
  name: string;
  shape: ValidVesselShapeType;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  capacity: number;
  heightInMillimeters: number;
  includeInGeneratedInstructions: boolean;
  widthInMillimeters: number;
}

export class ValidVesselCreationRequestInput implements IValidVesselCreationRequestInput {
  capacityUnitID?: string;
  displayInSummaryLists: boolean;
  lengthInMillimeters: number;
  description: string;
  iconPath: string;
  name: string;
  shape: ValidVesselShapeType;
  pluralName: string;
  slug: string;
  usableForStorage: boolean;
  capacity: number;
  heightInMillimeters: number;
  includeInGeneratedInstructions: boolean;
  widthInMillimeters: number;
  constructor(input: Partial<ValidVesselCreationRequestInput> = {}) {
    this.capacityUnitID = input.capacityUnitID;
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.lengthInMillimeters = input.lengthInMillimeters = 0;
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.name = input.name = '';
    this.shape = input.shape = 'other';
    this.pluralName = input.pluralName = '';
    this.slug = input.slug = '';
    this.usableForStorage = input.usableForStorage = false;
    this.capacity = input.capacity = 0;
    this.heightInMillimeters = input.heightInMillimeters = 0;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.widthInMillimeters = input.widthInMillimeters = 0;
  }
}
