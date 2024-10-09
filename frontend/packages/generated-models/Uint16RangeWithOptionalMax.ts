// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUint16RangeWithOptionalMax {
  max?: number;
  min: number;
}

export class Uint16RangeWithOptionalMax implements IUint16RangeWithOptionalMax {
  max?: number;
  min: number;
  constructor(input: Partial<Uint16RangeWithOptionalMax> = {}) {
    this.max = input.max;
    this.min = input.min = 0;
  }
}
