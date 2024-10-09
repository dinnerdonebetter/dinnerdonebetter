// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitCreationRequestInput {
  imperial: boolean;
  pluralName: string;
  universal: boolean;
  volumetric: boolean;
  description: string;
  metric: boolean;
  name: string;
  slug: string;
  iconPath: string;
}

export class ValidMeasurementUnitCreationRequestInput implements IValidMeasurementUnitCreationRequestInput {
  imperial: boolean;
  pluralName: string;
  universal: boolean;
  volumetric: boolean;
  description: string;
  metric: boolean;
  name: string;
  slug: string;
  iconPath: string;
  constructor(input: Partial<ValidMeasurementUnitCreationRequestInput> = {}) {
    this.imperial = input.imperial || false;
    this.pluralName = input.pluralName || '';
    this.universal = input.universal || false;
    this.volumetric = input.volumetric || false;
    this.description = input.description || '';
    this.metric = input.metric || false;
    this.name = input.name || '';
    this.slug = input.slug || '';
    this.iconPath = input.iconPath || '';
  }
}
