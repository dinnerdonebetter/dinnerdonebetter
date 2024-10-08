// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { OAuth2Client, APIResponse } from '@dinnerdonebetter/models';

export async function getOAuth2Client(client: Axios, oauth2ClientID: string): Promise<APIResponse<OAuth2Client>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<OAuth2Client>>(`/api/v1/oauth2_clients/${oauth2ClientID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
