// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUint16RangeWithOptionalMaxUpdateRequestInput {
   max?: number;
 min?: number;

}

export class Uint16RangeWithOptionalMaxUpdateRequestInput implements IUint16RangeWithOptionalMaxUpdateRequestInput {
   max?: number;
 min?: number;
constructor(input: Partial<Uint16RangeWithOptionalMaxUpdateRequestInput> = {}) {
	 this.max = input.max;
 this.min = input.min;
}
}