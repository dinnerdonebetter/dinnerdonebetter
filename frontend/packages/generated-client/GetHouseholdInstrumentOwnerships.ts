// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { HouseholdInstrumentOwnership, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getHouseholdInstrumentOwnerships(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<HouseholdInstrumentOwnership>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<HouseholdInstrumentOwnership>>>(
      `/api/v1/households/instruments`,
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<HouseholdInstrumentOwnership>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
