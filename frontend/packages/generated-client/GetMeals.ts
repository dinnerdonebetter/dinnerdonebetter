// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { Meal, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getMeals(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<Meal>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<Meal>>>(`/api/v1/meals`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Meal>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
