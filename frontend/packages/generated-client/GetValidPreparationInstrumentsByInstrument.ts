// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidPreparationInstrument, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getValidPreparationInstrumentsByInstrument(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  validInstrumentID: string,
): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ValidPreparationInstrument>>>(
      `/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`,
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidPreparationInstrument>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
