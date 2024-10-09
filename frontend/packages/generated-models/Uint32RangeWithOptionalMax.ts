// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUint32RangeWithOptionalMax {
  max?: number;
  min: number;
}

export class Uint32RangeWithOptionalMax implements IUint32RangeWithOptionalMax {
  max?: number;
  min: number;
  constructor(input: Partial<Uint32RangeWithOptionalMax> = {}) {
    this.max = input.max;
    this.min = input.min = 0;
  }
}
