// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidVessel, 
  APIResponse, 
  ValidVesselCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidVessel(
  client: Axios,
  
  input: ValidVesselCreationRequestInput,
): Promise<  APIResponse <  ValidVessel >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidVessel  >  >(`/api/v1/valid_vessels`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}