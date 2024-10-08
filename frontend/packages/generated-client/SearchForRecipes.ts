// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { Recipe, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function searchForRecipes(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  q: string,
): Promise<QueryFilteredResult<Recipe>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<Recipe>>>(`/api/v1/recipes/search`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Recipe>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
