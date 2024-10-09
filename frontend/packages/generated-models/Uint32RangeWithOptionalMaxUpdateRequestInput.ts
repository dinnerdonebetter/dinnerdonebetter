// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUint32RangeWithOptionalMaxUpdateRequestInput {
  max?: number;
  min?: number;
}

export class Uint32RangeWithOptionalMaxUpdateRequestInput implements IUint32RangeWithOptionalMaxUpdateRequestInput {
  max?: number;
  min?: number;
  constructor(input: Partial<Uint32RangeWithOptionalMaxUpdateRequestInput> = {}) {
    this.max = input.max;
    this.min = input.min;
  }
}
