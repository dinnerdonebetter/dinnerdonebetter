// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitUpdateRequestInput {
  description: string;
  iconPath: string;
  imperial: boolean;
  pluralName: string;
  volumetric: boolean;
  metric: boolean;
  name: string;
  slug: string;
  universal: boolean;
}

export class ValidMeasurementUnitUpdateRequestInput implements IValidMeasurementUnitUpdateRequestInput {
  description: string;
  iconPath: string;
  imperial: boolean;
  pluralName: string;
  volumetric: boolean;
  metric: boolean;
  name: string;
  slug: string;
  universal: boolean;
  constructor(input: Partial<ValidMeasurementUnitUpdateRequestInput> = {}) {
    this.description = input.description || '';
    this.iconPath = input.iconPath || '';
    this.imperial = input.imperial || false;
    this.pluralName = input.pluralName || '';
    this.volumetric = input.volumetric || false;
    this.metric = input.metric || false;
    this.name = input.name || '';
    this.slug = input.slug || '';
    this.universal = input.universal || false;
  }
}
