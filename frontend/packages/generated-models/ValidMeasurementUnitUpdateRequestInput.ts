// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IValidMeasurementUnitUpdateRequestInput {
  description?: string;
  iconPath?: string;
  name?: string;
  universal?: boolean;
  imperial?: boolean;
  metric?: boolean;
  pluralName?: string;
  slug?: string;
  volumetric?: boolean;
}

export class ValidMeasurementUnitUpdateRequestInput implements IValidMeasurementUnitUpdateRequestInput {
  description?: string;
  iconPath?: string;
  name?: string;
  universal?: boolean;
  imperial?: boolean;
  metric?: boolean;
  pluralName?: string;
  slug?: string;
  volumetric?: boolean;
  constructor(input: Partial<ValidMeasurementUnitUpdateRequestInput> = {}) {
    this.description = input.description;
    this.iconPath = input.iconPath;
    this.name = input.name;
    this.universal = input.universal;
    this.imperial = input.imperial;
    this.metric = input.metric;
    this.pluralName = input.pluralName;
    this.slug = input.slug;
    this.volumetric = input.volumetric;
  }
}
