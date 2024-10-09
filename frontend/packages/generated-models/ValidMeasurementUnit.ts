// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnit {
  lastUpdatedAt?: string;
  metric: boolean;
  name: string;
  archivedAt?: string;
  createdAt: string;
  description: string;
  iconPath: string;
  imperial: boolean;
  id: string;
  pluralName: string;
  slug: string;
  universal: boolean;
  volumetric: boolean;
}

export class ValidMeasurementUnit implements IValidMeasurementUnit {
  lastUpdatedAt?: string;
  metric: boolean;
  name: string;
  archivedAt?: string;
  createdAt: string;
  description: string;
  iconPath: string;
  imperial: boolean;
  id: string;
  pluralName: string;
  slug: string;
  universal: boolean;
  volumetric: boolean;
  constructor(input: Partial<ValidMeasurementUnit> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.metric = input.metric = false;
    this.name = input.name = '';
    this.archivedAt = input.archivedAt;
    this.createdAt = input.createdAt = '';
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.imperial = input.imperial = false;
    this.id = input.id = '';
    this.pluralName = input.pluralName = '';
    this.slug = input.slug = '';
    this.universal = input.universal = false;
    this.volumetric = input.volumetric = false;
  }
}
