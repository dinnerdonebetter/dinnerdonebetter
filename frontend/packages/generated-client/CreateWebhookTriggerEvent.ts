// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  WebhookTriggerEvent, 
  APIResponse, 
  WebhookTriggerEventCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createWebhookTriggerEvent(
  client: Axios,
  webhookID: string,
  input: WebhookTriggerEventCreationRequestInput,
): Promise<  APIResponse <  WebhookTriggerEvent >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < WebhookTriggerEvent  >  >(`/api/v1/webhooks/${webhookID}/trigger_events`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}