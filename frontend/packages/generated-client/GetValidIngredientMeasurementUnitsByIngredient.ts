// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientMeasurementUnit,
  QueryFilter,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

export async function getValidIngredientMeasurementUnitsByIngredient(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  validIngredientID: string,
): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ValidIngredientMeasurementUnit>>>(
      `/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`,
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientMeasurementUnit>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
