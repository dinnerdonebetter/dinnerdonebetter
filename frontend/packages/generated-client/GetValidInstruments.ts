// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidInstrument, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getValidInstruments(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ValidInstrument>>>(`/api/v1/valid_instruments`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidInstrument>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
