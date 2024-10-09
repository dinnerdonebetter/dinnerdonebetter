// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidPreparationVessel, 
  APIResponse, 
  ValidPreparationVesselCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidPreparationVessel(
  client: Axios,
  
  input: ValidPreparationVesselCreationRequestInput,
): Promise<  APIResponse <  ValidPreparationVessel >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidPreparationVessel  >  >(`/api/v1/valid_preparation_vessels`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}