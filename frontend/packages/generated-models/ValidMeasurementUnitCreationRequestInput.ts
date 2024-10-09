// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitCreationRequestInput {
  name: string;
  volumetric: boolean;
  description: string;
  iconPath: string;
  imperial: boolean;
  metric: boolean;
  pluralName: string;
  slug: string;
  universal: boolean;
}

export class ValidMeasurementUnitCreationRequestInput implements IValidMeasurementUnitCreationRequestInput {
  name: string;
  volumetric: boolean;
  description: string;
  iconPath: string;
  imperial: boolean;
  metric: boolean;
  pluralName: string;
  slug: string;
  universal: boolean;
  constructor(input: Partial<ValidMeasurementUnitCreationRequestInput> = {}) {
    this.name = input.name = '';
    this.volumetric = input.volumetric = false;
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.imperial = input.imperial = false;
    this.metric = input.metric = false;
    this.pluralName = input.pluralName = '';
    this.slug = input.slug = '';
    this.universal = input.universal = false;
  }
}
