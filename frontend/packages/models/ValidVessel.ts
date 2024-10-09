// GENERATED CODE, DO NOT EDIT MANUALLY

import { ValidMeasurementUnit } from './ValidMeasurementUnit';
import { ValidVesselShapeType } from './enums';

export interface IValidVessel {
  widthInMillimeters: number;
  archivedAt: string;
  capacity: number;
  capacityUnit: ValidMeasurementUnit;
  heightInMillimeters: number;
  name: string;
  slug: string;
  usableForStorage: boolean;
  createdAt: string;
  description: string;
  displayInSummaryLists: boolean;
  id: string;
  shape: ValidVesselShapeType;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt: string;
  lengthInMillimeters: number;
  pluralName: string;
}

export class ValidVessel implements IValidVessel {
  widthInMillimeters: number;
  archivedAt: string;
  capacity: number;
  capacityUnit: ValidMeasurementUnit;
  heightInMillimeters: number;
  name: string;
  slug: string;
  usableForStorage: boolean;
  createdAt: string;
  description: string;
  displayInSummaryLists: boolean;
  id: string;
  shape: ValidVesselShapeType;
  iconPath: string;
  includeInGeneratedInstructions: boolean;
  lastUpdatedAt: string;
  lengthInMillimeters: number;
  pluralName: string;
  constructor(input: Partial<ValidVessel> = {}) {
    this.widthInMillimeters = input.widthInMillimeters || 0;
    this.archivedAt = input.archivedAt || '';
    this.capacity = input.capacity || 0;
    this.capacityUnit = input.capacityUnit || new ValidMeasurementUnit();
    this.heightInMillimeters = input.heightInMillimeters || 0;
    this.name = input.name || '';
    this.slug = input.slug || '';
    this.usableForStorage = input.usableForStorage || false;
    this.createdAt = input.createdAt || '';
    this.description = input.description || '';
    this.displayInSummaryLists = input.displayInSummaryLists || false;
    this.id = input.id || '';
    this.shape = input.shape || 'other';
    this.iconPath = input.iconPath || '';
    this.includeInGeneratedInstructions = input.includeInGeneratedInstructions || false;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.lengthInMillimeters = input.lengthInMillimeters || 0;
    this.pluralName = input.pluralName || '';
  }
}
