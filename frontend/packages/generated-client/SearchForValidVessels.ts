// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidVessel, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function searchForValidVessels(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  q: string,
): Promise<QueryFilteredResult<ValidVessel>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ValidVessel>>>(`/api/v1/valid_vessels/search`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidVessel>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
