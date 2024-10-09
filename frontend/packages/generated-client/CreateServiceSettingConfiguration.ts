// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ServiceSettingConfiguration, 
  APIResponse, 
  ServiceSettingConfigurationCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createServiceSettingConfiguration(
  client: Axios,
  
  input: ServiceSettingConfigurationCreationRequestInput,
): Promise<  APIResponse <  ServiceSettingConfiguration >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ServiceSettingConfiguration  >  >(`/api/v1/settings/configurations`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}