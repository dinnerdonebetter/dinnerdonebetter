// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidVesselShapeType } from './enums';

export interface IValidVessel {
  heightInMillimeters: number;
  pluralName: string;
  capacityUnit?: ValidMeasurementUnit;
  iconPath: string;
  id: string;
  lengthInMillimeters: number;
  slug: string;
  widthInMillimeters: number;
  archivedAt?: string;
  createdAt: string;
  description: string;
  name: string;
  shape: ValidVesselShapeType;
  capacity: number;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt?: string;
  usableForStorage: boolean;
  displayInSummaryLists: boolean;
}

export class ValidVessel implements IValidVessel {
  heightInMillimeters: number;
  pluralName: string;
  capacityUnit?: ValidMeasurementUnit;
  iconPath: string;
  id: string;
  lengthInMillimeters: number;
  slug: string;
  widthInMillimeters: number;
  archivedAt?: string;
  createdAt: string;
  description: string;
  name: string;
  shape: ValidVesselShapeType;
  capacity: number;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt?: string;
  usableForStorage: boolean;
  displayInSummaryLists: boolean;
  constructor(input: Partial<ValidVessel> = {}) {
    this.heightInMillimeters = input.heightInMillimeters = 0;
    this.pluralName = input.pluralName = '';
    this.capacityUnit = input.capacityUnit;
    this.iconPath = input.iconPath = '';
    this.id = input.id = '';
    this.lengthInMillimeters = input.lengthInMillimeters = 0;
    this.slug = input.slug = '';
    this.widthInMillimeters = input.widthInMillimeters = 0;
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.name = input.name = '';
    this.shape = input.shape = 'other';
    this.capacity = input.capacity = 0;
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.usableForStorage = input.usableForStorage = false;
    this.displayInSummaryLists = input.displayInSummaryLists = false;
  }
}
