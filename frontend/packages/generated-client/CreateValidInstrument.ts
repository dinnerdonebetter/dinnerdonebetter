// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidInstrument, 
  APIResponse, 
  ValidInstrumentCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidInstrument(
  client: Axios,
  
  input: ValidInstrumentCreationRequestInput,
): Promise<  APIResponse <  ValidInstrument >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidInstrument  >  >(`/api/v1/valid_instruments`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}