// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ServiceSettingConfiguration, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveServiceSettingConfiguration(
  client: Axios,
  serviceSettingConfigurationID: string,
	): Promise<  APIResponse <  ServiceSettingConfiguration >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < ServiceSettingConfiguration  >  >(`/api/v1/settings/configurations/${serviceSettingConfigurationID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}