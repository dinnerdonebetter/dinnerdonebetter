// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IFloat32RangeWithOptionalMax {
  min: number;
  max?: number;
}

export class Float32RangeWithOptionalMax implements IFloat32RangeWithOptionalMax {
  min: number;
  max?: number;
  constructor(input: Partial<Float32RangeWithOptionalMax> = {}) {
    this.min = input.min = 0;
    this.max = input.max;
  }
}
