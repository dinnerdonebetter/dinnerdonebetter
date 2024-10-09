// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  WebhookTriggerEvent, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveWebhookTriggerEvent(
  client: Axios,
  webhookID: string,
	webhookTriggerEventID: string,
	): Promise<  APIResponse <  WebhookTriggerEvent >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < WebhookTriggerEvent  >  >(`/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}