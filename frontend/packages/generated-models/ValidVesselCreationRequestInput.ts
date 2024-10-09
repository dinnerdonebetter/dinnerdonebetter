// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVesselShapeType } from './enums';

export interface IValidVesselCreationRequestInput {
  widthInMillimeters: number;
  lengthInMillimeters: number;
  name: string;
  pluralName: string;
  shape: ValidVesselShapeType;
  usableForStorage: boolean;
  description: string;
  heightInMillimeters: number;
  iconPath: string;
  slug: string;
  capacityUnitID?: string;
  capacity: number;
  displayInSummaryLists: boolean;
  includeInGeneratedInstructions: boolean;
}

export class ValidVesselCreationRequestInput implements IValidVesselCreationRequestInput {
  widthInMillimeters: number;
  lengthInMillimeters: number;
  name: string;
  pluralName: string;
  shape: ValidVesselShapeType;
  usableForStorage: boolean;
  description: string;
  heightInMillimeters: number;
  iconPath: string;
  slug: string;
  capacityUnitID?: string;
  capacity: number;
  displayInSummaryLists: boolean;
  includeInGeneratedInstructions: boolean;
  constructor(input: Partial<ValidVesselCreationRequestInput> = {}) {
    this.widthInMillimeters = input.widthInMillimeters = 0;
    this.lengthInMillimeters = input.lengthInMillimeters = 0;
    this.name = input.name = '';
    this.pluralName = input.pluralName = '';
    this.shape = input.shape = 'other';
    this.usableForStorage = input.usableForStorage = false;
    this.description = input.description = '';
    this.heightInMillimeters = input.heightInMillimeters = 0;
    this.iconPath = input.iconPath = '';
    this.slug = input.slug = '';
    this.capacityUnitID = input.capacityUnitID;
    this.capacity = input.capacity = 0;
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
  }
}
