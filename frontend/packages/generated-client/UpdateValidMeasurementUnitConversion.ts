// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidMeasurementUnitConversion, 
  APIResponse, 
  ValidMeasurementUnitConversionUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateValidMeasurementUnitConversion(
  client: Axios,
  validMeasurementUnitConversionID: string,
  input: ValidMeasurementUnitConversionUpdateRequestInput,
): Promise<  APIResponse <  ValidMeasurementUnitConversion >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < ValidMeasurementUnitConversion  >  >(`/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}