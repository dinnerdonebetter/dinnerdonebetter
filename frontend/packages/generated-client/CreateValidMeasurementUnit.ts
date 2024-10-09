// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidMeasurementUnit, 
  APIResponse, 
  ValidMeasurementUnitCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidMeasurementUnit(
  client: Axios,
  
  input: ValidMeasurementUnitCreationRequestInput,
): Promise<  APIResponse <  ValidMeasurementUnit >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidMeasurementUnit  >  >(`/api/v1/valid_measurement_units`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}