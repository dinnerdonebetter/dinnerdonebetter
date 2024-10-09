// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitUpdateRequestInput {
  universal?: boolean;
  volumetric?: boolean;
  name?: string;
  iconPath?: string;
  imperial?: boolean;
  metric?: boolean;
  pluralName?: string;
  slug?: string;
  description?: string;
}

export class ValidMeasurementUnitUpdateRequestInput implements IValidMeasurementUnitUpdateRequestInput {
  universal?: boolean;
  volumetric?: boolean;
  name?: string;
  iconPath?: string;
  imperial?: boolean;
  metric?: boolean;
  pluralName?: string;
  slug?: string;
  description?: string;
  constructor(input: Partial<ValidMeasurementUnitUpdateRequestInput> = {}) {
    this.universal = input.universal;
    this.volumetric = input.volumetric;
    this.name = input.name;
    this.iconPath = input.iconPath;
    this.imperial = input.imperial;
    this.metric = input.metric;
    this.pluralName = input.pluralName;
    this.slug = input.slug;
    this.description = input.description;
  }
}
