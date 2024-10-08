// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { Webhook, APIResponse } from '@dinnerdonebetter/models';

export async function getWebhook(client: Axios, webhookID: string): Promise<APIResponse<Webhook>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Webhook>>(`/api/v1/webhooks/${webhookID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
