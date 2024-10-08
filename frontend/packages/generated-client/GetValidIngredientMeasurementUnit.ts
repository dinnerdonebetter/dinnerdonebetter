// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidIngredientMeasurementUnit, APIResponse } from '@dinnerdonebetter/models';

export async function getValidIngredientMeasurementUnit(
  client: Axios,
  validIngredientMeasurementUnitID: string,
): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredientMeasurementUnit>>(
      `/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
