// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUint32RangeWithOptionalMaxUpdateRequestInput {
   min?: number;
 max?: number;

}

export class Uint32RangeWithOptionalMaxUpdateRequestInput implements IUint32RangeWithOptionalMaxUpdateRequestInput {
   min?: number;
 max?: number;
constructor(input: Partial<Uint32RangeWithOptionalMaxUpdateRequestInput> = {}) {
	 this.min = input.min;
 this.max = input.max;
}
}