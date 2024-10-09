// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ServiceSetting, 
  APIResponse, 
  ServiceSettingCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createServiceSetting(
  client: Axios,
  
  input: ServiceSettingCreationRequestInput,
): Promise<  APIResponse <  ServiceSetting >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ServiceSetting  >  >(`/api/v1/settings`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}