// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidVessel, 
  APIResponse, 
  ValidVesselUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateValidVessel(
  client: Axios,
  validVesselID: string,
  input: ValidVesselUpdateRequestInput,
): Promise<  APIResponse <  ValidVessel >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < ValidVessel  >  >(`/api/v1/valid_vessels/${validVesselID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}