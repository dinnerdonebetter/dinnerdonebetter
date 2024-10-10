// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnit {
  archivedAt: string;
  createdAt: string;
  description: string;
  iconPath: string;
  id: string;
  imperial: boolean;
  lastUpdatedAt: string;
  metric: boolean;
  name: string;
  pluralName: string;
  slug: string;
  universal: boolean;
  volumetric: boolean;
}

export class ValidMeasurementUnit implements IValidMeasurementUnit {
  archivedAt: string;
  createdAt: string;
  description: string;
  iconPath: string;
  id: string;
  imperial: boolean;
  lastUpdatedAt: string;
  metric: boolean;
  name: string;
  pluralName: string;
  slug: string;
  universal: boolean;
  volumetric: boolean;
  constructor(input: Partial<ValidMeasurementUnit> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
    this.id = input.id || '';
    this.imperial = input.imperial || false;
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.metric = input.metric || false;
    this.name = input.name || '';
    this.pluralName = input.pluralName || '';
    this.slug = input.slug || '';
    this.universal = input.universal || false;
    this.volumetric = input.volumetric || false;
  }
}
