// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ServiceSettingConfiguration,
  APIResponse,
  ServiceSettingConfigurationUpdateRequestInput,
} from '@dinnerdonebetter/models';

export async function updateServiceSettingConfiguration(
  client: Axios,
  serviceSettingConfigurationID: string,
  input: ServiceSettingConfigurationUpdateRequestInput,
): Promise<APIResponse<ServiceSettingConfiguration>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ServiceSettingConfiguration>>(
      `/api/v1/settings/configurations/${serviceSettingConfigurationID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
