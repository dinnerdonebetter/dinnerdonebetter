// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidMeasurementUnitConversion, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getValidMeasurementUnitConversionsFromUnit(
  client: Axios,
  validMeasurementUnitID: string,
	): Promise<  APIResponse <  ValidMeasurementUnitConversion >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < ValidMeasurementUnitConversion  >  >(`/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}