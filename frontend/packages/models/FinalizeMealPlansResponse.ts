// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IFinalizeMealPlansResponse {
   count: number;

}

export class FinalizeMealPlansResponse implements IFinalizeMealPlansResponse {
   count: number;
constructor(input: Partial<FinalizeMealPlansResponse> = {}) {
	 this.count = input.count || 0;
}
}