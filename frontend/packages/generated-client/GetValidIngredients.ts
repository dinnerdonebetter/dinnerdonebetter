// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidIngredient, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getValidIngredients(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidIngredient>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ValidIngredient>>>(`/api/v1/valid_ingredients`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredient>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
