// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IAPIError {
   code: string;
 message: string;

}

export class APIError implements IAPIError {
   code: string;
 message: string;
constructor(input: Partial<APIError> = {}) {
	 this.code = input.code || '';
 this.message = input.message || '';
}
}