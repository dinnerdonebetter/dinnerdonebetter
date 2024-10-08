// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidMeasurementUnit, APIResponse } from '@dinnerdonebetter/models';

export async function getValidMeasurementUnit(
  client: Axios,
  validMeasurementUnitID: string,
): Promise<APIResponse<ValidMeasurementUnit>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidMeasurementUnit>>(
      `/api/v1/valid_measurement_units/${validMeasurementUnitID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
