// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { OAuth2ClientCreationResponse, APIResponse, OAuth2ClientCreationRequestInput } from '@dinnerdonebetter/models';

export async function createOAuth2Client(
  client: Axios,

  input: OAuth2ClientCreationRequestInput,
): Promise<APIResponse<OAuth2ClientCreationResponse>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<OAuth2ClientCreationResponse>>(`/api/v1/oauth2_clients`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
