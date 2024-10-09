// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  Webhook, 
  APIResponse, 
  WebhookCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createWebhook(
  client: Axios,
  
  input: WebhookCreationRequestInput,
): Promise<  APIResponse <  Webhook >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < Webhook  >  >(`/api/v1/webhooks`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}