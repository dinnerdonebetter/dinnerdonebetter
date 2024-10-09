// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUint32RangeWithOptionalMax {
   min: number;
 max?: number;

}

export class Uint32RangeWithOptionalMax implements IUint32RangeWithOptionalMax {
   min: number;
 max?: number;
constructor(input: Partial<Uint32RangeWithOptionalMax> = {}) {
	 this.min = input.min = 0;
 this.max = input.max;
}
}