// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ServiceSetting, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getServiceSettings(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ServiceSetting>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ServiceSetting>>>(`/api/v1/settings`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ServiceSetting>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
