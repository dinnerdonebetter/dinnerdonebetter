// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ServiceSettingConfiguration, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getServiceSettingConfigurationByName(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  serviceSettingConfigurationName: string,
): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<ServiceSettingConfiguration>>>(
      `/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`,
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ServiceSettingConfiguration>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
