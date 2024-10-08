// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientMeasurementUnit,
  APIResponse,
  ValidIngredientMeasurementUnitCreationRequestInput,
} from '@dinnerdonebetter/models';

export async function createValidIngredientMeasurementUnit(
  client: Axios,

  input: ValidIngredientMeasurementUnitCreationRequestInput,
): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidIngredientMeasurementUnit>>(
      `/api/v1/valid_ingredient_measurement_units`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
