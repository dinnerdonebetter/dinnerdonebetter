// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IFloat32RangeWithOptionalMaxUpdateRequestInput {
   max: number;
 min: number;

}

export class Float32RangeWithOptionalMaxUpdateRequestInput implements IFloat32RangeWithOptionalMaxUpdateRequestInput {
   max: number;
 min: number;
constructor(input: Partial<Float32RangeWithOptionalMaxUpdateRequestInput> = {}) {
	 this.max = input.max || 0;
 this.min = input.min || 0;
}
}