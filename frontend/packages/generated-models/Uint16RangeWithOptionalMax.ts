// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUint16RangeWithOptionalMax {
  min: number;
  max?: number;
}

export class Uint16RangeWithOptionalMax implements IUint16RangeWithOptionalMax {
  min: number;
  max?: number;
  constructor(input: Partial<Uint16RangeWithOptionalMax> = {}) {
    this.min = input.min = 0;
    this.max = input.max;
  }
}
