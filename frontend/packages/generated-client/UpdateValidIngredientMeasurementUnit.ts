// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientMeasurementUnit,
  APIResponse,
  ValidIngredientMeasurementUnitUpdateRequestInput,
} from '@dinnerdonebetter/models';

export async function updateValidIngredientMeasurementUnit(
  client: Axios,
  validIngredientMeasurementUnitID: string,
  input: ValidIngredientMeasurementUnitUpdateRequestInput,
): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidIngredientMeasurementUnit>>(
      `/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
