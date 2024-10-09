// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnit {
  iconPath: string;
  imperial: boolean;
  name: string;
  pluralName: string;
  slug: string;
  universal: boolean;
  volumetric: boolean;
  createdAt: string;
  description: string;
  id: string;
  lastUpdatedAt?: string;
  metric: boolean;
  archivedAt?: string;
}

export class ValidMeasurementUnit implements IValidMeasurementUnit {
  iconPath: string;
  imperial: boolean;
  name: string;
  pluralName: string;
  slug: string;
  universal: boolean;
  volumetric: boolean;
  createdAt: string;
  description: string;
  id: string;
  lastUpdatedAt?: string;
  metric: boolean;
  archivedAt?: string;
  constructor(input: Partial<ValidMeasurementUnit> = {}) {
    this.iconPath = input.iconPath = '';
    this.imperial = input.imperial = false;
    this.name = input.name = '';
    this.pluralName = input.pluralName = '';
    this.slug = input.slug = '';
    this.universal = input.universal = false;
    this.volumetric = input.volumetric = false;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.metric = input.metric = false;
    this.archivedAt = input.archivedAt;
  }
}
