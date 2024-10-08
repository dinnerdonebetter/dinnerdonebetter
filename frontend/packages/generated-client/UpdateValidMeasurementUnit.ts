// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidMeasurementUnit, APIResponse, ValidMeasurementUnitUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateValidMeasurementUnit(
  client: Axios,
  validMeasurementUnitID: string,
  input: ValidMeasurementUnitUpdateRequestInput,
): Promise<APIResponse<ValidMeasurementUnit>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidMeasurementUnit>>(
      `/api/v1/valid_measurement_units/${validMeasurementUnitID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
