// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IOptionalFloat32Range {
   max: number;
 min: number;

}

export class OptionalFloat32Range implements IOptionalFloat32Range {
   max: number;
 min: number;
constructor(input: Partial<OptionalFloat32Range> = {}) {
	 this.max = input.max || 0;
 this.min = input.min || 0;
}
}