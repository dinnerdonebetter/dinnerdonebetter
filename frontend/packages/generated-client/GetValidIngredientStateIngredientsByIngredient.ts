// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientStateIngredient,
  QueryFilter,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

export async function getValidIngredientStateIngredientsByIngredient(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  validIngredientID: string,
): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ValidIngredientStateIngredient>>>(
      `/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`,
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientStateIngredient>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
