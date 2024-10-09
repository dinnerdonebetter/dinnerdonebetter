// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ServiceSetting, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveServiceSetting(
  client: Axios,
  serviceSettingID: string,
	): Promise<  APIResponse <  ServiceSetting >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < ServiceSetting  >  >(`/api/v1/settings/${serviceSettingID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}