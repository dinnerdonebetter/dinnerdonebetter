// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidPreparation, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function searchForValidPreparations(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  q: string,
): Promise<QueryFilteredResult<ValidPreparation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ValidPreparation>>>(`/api/v1/valid_preparations/search`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidPreparation>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
