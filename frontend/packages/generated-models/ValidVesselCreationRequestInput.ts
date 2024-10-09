// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidVesselShapeType } from './enums';

export interface IValidVesselCreationRequestInput {
  includeInGeneratedInstructions: boolean;
  widthInMillimeters: number;
  iconPath: string;
  pluralName: string;
  capacityUnitID?: string;
  displayInSummaryLists: boolean;
  lengthInMillimeters: number;
  name: string;
  shape: ValidVesselShapeType;
  slug: string;
  usableForStorage: boolean;
  capacity: number;
  description: string;
  heightInMillimeters: number;
}

export class ValidVesselCreationRequestInput implements IValidVesselCreationRequestInput {
  includeInGeneratedInstructions: boolean;
  widthInMillimeters: number;
  iconPath: string;
  pluralName: string;
  capacityUnitID?: string;
  displayInSummaryLists: boolean;
  lengthInMillimeters: number;
  name: string;
  shape: ValidVesselShapeType;
  slug: string;
  usableForStorage: boolean;
  capacity: number;
  description: string;
  heightInMillimeters: number;
  constructor(input: Partial<ValidVesselCreationRequestInput> = {}) {
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.widthInMillimeters = input.widthInMillimeters = 0;
    this.iconPath = input.iconPath = '';
    this.pluralName = input.pluralName = '';
    this.capacityUnitID = input.capacityUnitID;
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.lengthInMillimeters = input.lengthInMillimeters = 0;
    this.name = input.name = '';
    this.shape = input.shape = 'other';
    this.slug = input.slug = '';
    this.usableForStorage = input.usableForStorage = false;
    this.capacity = input.capacity = 0;
    this.description = input.description = '';
    this.heightInMillimeters = input.heightInMillimeters = 0;
  }
}
