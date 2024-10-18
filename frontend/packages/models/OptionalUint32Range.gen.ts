// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IOptionalUint32Range {
  max: number;
  min: number;
}

export class OptionalUint32Range implements IOptionalUint32Range {
  max: number;
  min: number;
  constructor(input: Partial<OptionalUint32Range> = {}) {
    this.max = input.max || 0;
    this.min = input.min || 0;
  }
}
