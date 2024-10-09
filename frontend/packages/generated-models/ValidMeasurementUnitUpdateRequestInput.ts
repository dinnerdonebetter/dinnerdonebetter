// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitUpdateRequestInput {
  imperial?: boolean;
  name?: string;
  pluralName?: string;
  slug?: string;
  description?: string;
  metric?: boolean;
  universal?: boolean;
  volumetric?: boolean;
  iconPath?: string;
}

export class ValidMeasurementUnitUpdateRequestInput implements IValidMeasurementUnitUpdateRequestInput {
  imperial?: boolean;
  name?: string;
  pluralName?: string;
  slug?: string;
  description?: string;
  metric?: boolean;
  universal?: boolean;
  volumetric?: boolean;
  iconPath?: string;
  constructor(input: Partial<ValidMeasurementUnitUpdateRequestInput> = {}) {
    this.imperial = input.imperial;
    this.name = input.name;
    this.pluralName = input.pluralName;
    this.slug = input.slug;
    this.description = input.description;
    this.metric = input.metric;
    this.universal = input.universal;
    this.volumetric = input.volumetric;
    this.iconPath = input.iconPath;
  }
}
