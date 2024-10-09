// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidPreparationVessel, 
  APIResponse, 
  ValidPreparationVesselUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateValidPreparationVessel(
  client: Axios,
  validPreparationVesselID: string,
  input: ValidPreparationVesselUpdateRequestInput,
): Promise<  APIResponse <  ValidPreparationVessel >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < ValidPreparationVessel  >  >(`/api/v1/valid_preparation_vessels/${validPreparationVesselID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}