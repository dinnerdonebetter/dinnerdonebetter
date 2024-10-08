// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidMeasurementUnitConversion, APIResponse } from '@dinnerdonebetter/models';

export async function getValidMeasurementUnitConversion(
  client: Axios,
  validMeasurementUnitConversionID: string,
): Promise<APIResponse<ValidMeasurementUnitConversion>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidMeasurementUnitConversion>>(
      `/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
