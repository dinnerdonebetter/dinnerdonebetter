// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnit {
  imperial: boolean;
  pluralName: string;
  slug: string;
  universal: boolean;
  volumetric: boolean;
  archivedAt: string;
  description: string;
  iconPath: string;
  metric: boolean;
  name: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
}

export class ValidMeasurementUnit implements IValidMeasurementUnit {
  imperial: boolean;
  pluralName: string;
  slug: string;
  universal: boolean;
  volumetric: boolean;
  archivedAt: string;
  description: string;
  iconPath: string;
  metric: boolean;
  name: string;
  createdAt: string;
  id: string;
  lastUpdatedAt: string;
  constructor(input: Partial<ValidMeasurementUnit> = {}) {
    this.imperial = input.imperial || false;
    this.pluralName = input.pluralName || '';
    this.slug = input.slug || '';
    this.universal = input.universal || false;
    this.volumetric = input.volumetric || false;
    this.archivedAt = input.archivedAt || '';
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
    this.metric = input.metric || false;
    this.name = input.name || '';
    this.createdAt = input.createdAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
  }
}
