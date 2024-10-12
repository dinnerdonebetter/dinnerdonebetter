// GENERATED CODE, DO NOT EDIT MANUALLY

 import { APIError } from './APIError';
 import { ResponseDetails } from './ResponseDetails';


export interface IAPIResponseWithError {
   details: ResponseDetails;
 error: APIError;

}

export class APIResponseWithError implements IAPIResponseWithError {
   details: ResponseDetails;
 error: APIError;
constructor(input: Partial<APIResponseWithError> = {}) {
	 this.details = input.details || new ResponseDetails();
 this.error = input.error || new APIError();
}
}