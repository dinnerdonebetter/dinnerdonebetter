// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnit {
  description: string;
  imperial: boolean;
  pluralName: string;
  universal: boolean;
  volumetric: boolean;
  slug: string;
  archivedAt?: string;
  createdAt: string;
  iconPath: string;
  id: string;
  lastUpdatedAt?: string;
  metric: boolean;
  name: string;
}

export class ValidMeasurementUnit implements IValidMeasurementUnit {
  description: string;
  imperial: boolean;
  pluralName: string;
  universal: boolean;
  volumetric: boolean;
  slug: string;
  archivedAt?: string;
  createdAt: string;
  iconPath: string;
  id: string;
  lastUpdatedAt?: string;
  metric: boolean;
  name: string;
  constructor(input: Partial<ValidMeasurementUnit> = {}) {
    this.description = input.description = '';
    this.imperial = input.imperial = false;
    this.pluralName = input.pluralName = '';
    this.universal = input.universal = false;
    this.volumetric = input.volumetric = false;
    this.slug = input.slug = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.iconPath = input.iconPath = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.metric = input.metric = false;
    this.name = input.name = '';
  }
}
