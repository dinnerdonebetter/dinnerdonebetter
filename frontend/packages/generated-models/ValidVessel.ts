// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidVesselShapeType } from './enums';

export interface IValidVessel {
  capacityUnit?: ValidMeasurementUnit;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  shape: ValidVesselShapeType;
  usableForStorage: boolean;
  lastUpdatedAt?: string;
  pluralName: string;
  widthInMillimeters: number;
  id: string;
  archivedAt?: string;
  capacity: number;
  createdAt: string;
  heightInMillimeters: number;
  lengthInMillimeters: number;
  name: string;
  slug: string;
}

export class ValidVessel implements IValidVessel {
  capacityUnit?: ValidMeasurementUnit;
  description: string;
  displayInSummaryLists: boolean;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  shape: ValidVesselShapeType;
  usableForStorage: boolean;
  lastUpdatedAt?: string;
  pluralName: string;
  widthInMillimeters: number;
  id: string;
  archivedAt?: string;
  capacity: number;
  createdAt: string;
  heightInMillimeters: number;
  lengthInMillimeters: number;
  name: string;
  slug: string;
  constructor(input: Partial<ValidVessel> = {}) {
    this.capacityUnit = input.capacityUnit;
    this.description = input.description = '';
    this.displayInSummaryLists = input.displayInSummaryLists = false;
    this.iconPath = input.iconPath = '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions = false;
    this.shape = input.shape = 'other';
    this.usableForStorage = input.usableForStorage = false;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.pluralName = input.pluralName = '';
    this.widthInMillimeters = input.widthInMillimeters = 0;
    this.id = input.id = '';
    this.archivedAt = input.archivedAt;
    this.capacity = input.capacity = 0;
    this.createdAt = input.createdAt = '';
    this.heightInMillimeters = input.heightInMillimeters = 0;
    this.lengthInMillimeters = input.lengthInMillimeters = 0;
    this.name = input.name = '';
    this.slug = input.slug = '';
  }
}
