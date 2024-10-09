// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitCreationRequestInput {
  volumetric: boolean;
  slug: string;
  universal: boolean;
  description: string;
  iconPath: string;
  imperial: boolean;
  metric: boolean;
  name: string;
  pluralName: string;
}

export class ValidMeasurementUnitCreationRequestInput implements IValidMeasurementUnitCreationRequestInput {
  volumetric: boolean;
  slug: string;
  universal: boolean;
  description: string;
  iconPath: string;
  imperial: boolean;
  metric: boolean;
  name: string;
  pluralName: string;
  constructor(input: Partial<ValidMeasurementUnitCreationRequestInput> = {}) {
    this.volumetric = input.volumetric = false;
    this.slug = input.slug = '';
    this.universal = input.universal = false;
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.imperial = input.imperial = false;
    this.metric = input.metric = false;
    this.name = input.name = '';
    this.pluralName = input.pluralName = '';
  }
}
