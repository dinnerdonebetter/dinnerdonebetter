// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitCreationRequestInput {
  metric: boolean;
  slug: string;
  volumetric: boolean;
  pluralName: string;
  universal: boolean;
  description: string;
  iconPath: string;
  imperial: boolean;
  name: string;
}

export class ValidMeasurementUnitCreationRequestInput implements IValidMeasurementUnitCreationRequestInput {
  metric: boolean;
  slug: string;
  volumetric: boolean;
  pluralName: string;
  universal: boolean;
  description: string;
  iconPath: string;
  imperial: boolean;
  name: string;
  constructor(input: Partial<ValidMeasurementUnitCreationRequestInput> = {}) {
    this.metric = input.metric = false;
    this.slug = input.slug = '';
    this.volumetric = input.volumetric = false;
    this.pluralName = input.pluralName = '';
    this.universal = input.universal = false;
    this.description = input.description = '';
    this.iconPath = input.iconPath = '';
    this.imperial = input.imperial = false;
    this.name = input.name = '';
  }
}
