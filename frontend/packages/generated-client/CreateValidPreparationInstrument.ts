// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidPreparationInstrument, 
  APIResponse, 
  ValidPreparationInstrumentCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidPreparationInstrument(
  client: Axios,
  
  input: ValidPreparationInstrumentCreationRequestInput,
): Promise<  APIResponse <  ValidPreparationInstrument >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidPreparationInstrument  >  >(`/api/v1/valid_preparation_instruments`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}