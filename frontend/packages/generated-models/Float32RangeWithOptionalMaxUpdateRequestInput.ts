// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IFloat32RangeWithOptionalMaxUpdateRequestInput {
  min?: number;
  max?: number;
}

export class Float32RangeWithOptionalMaxUpdateRequestInput implements IFloat32RangeWithOptionalMaxUpdateRequestInput {
  min?: number;
  max?: number;
  constructor(input: Partial<Float32RangeWithOptionalMaxUpdateRequestInput> = {}) {
    this.min = input.min;
    this.max = input.max;
  }
}
