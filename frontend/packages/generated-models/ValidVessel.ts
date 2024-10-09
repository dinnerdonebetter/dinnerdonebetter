// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidVesselShapeType } from './enums';

export interface IValidVessel {
  capacity: number;
  createdAt: string;
  iconPath: string;
  name: string;
  slug: string;
  id: string;
  includeInGeneratedInstructions: boolean;
  lengthInMillimeters: number;
  pluralName: string;
  usableForStorage: boolean;
  capacityUnit?: ValidMeasurementUnit;
  displayInSummaryLists: boolean;
  heightInMillimeters: number;
  lastUpdatedAt?: string;
  archivedAt?: string;
  description: string;
  shape: ValidVesselShapeType;
  widthInMillimeters: number;
}

export class ValidVessel implements IValidVessel {
  capacity: number;
  createdAt: string;
  iconPath: string;
  name: string;
  slug: string;
  id: string;
  includeInGeneratedInstructions: boolean;
  lengthInMillimeters: number;
  pluralName: string;
  usableForStorage: boolean;
  capacityUnit?: ValidMeasurementUnit;
  displayInSummaryLists: boolean;
  heightInMillimeters: number;
  lastUpdatedAt?: string;
  archivedAt?: string;
  description: string;
  shape: ValidVesselShapeType;
  widthInMillimeters: number;
  constructor(input: Partial<ValidVessel> = {}) {
    this.capacity = input.capacity = 0;
    this.createdAt = input.createdAt = '';
    this.iconPath = input.iconPath = '';
    this.name = input.name = '';
    this.slug = input.slug = '';
    this.id = input.id = '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.lengthInMillimeters = input.lengthInMillimeters = 0;
    this.pluralName = input.pluralName = '';
    this.usableForStorage = input.usableForStorage = false;
    this.capacityUnit = input.capacityUnit;
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.heightInMillimeters = input.heightInMillimeters = 0;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.archivedAt = input.archivedAt;
    this.description = input.description = '';
    this.shape = input.shape = 'other';
    this.widthInMillimeters = input.widthInMillimeters = 0;
  }
}
