// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IFloat32RangeWithOptionalMax {
   max?: number;
 min: number;

}

export class Float32RangeWithOptionalMax implements IFloat32RangeWithOptionalMax {
   max?: number;
 min: number;
constructor(input: Partial<Float32RangeWithOptionalMax> = {}) {
	 this.max = input.max;
 this.min = input.min = 0;
}
}