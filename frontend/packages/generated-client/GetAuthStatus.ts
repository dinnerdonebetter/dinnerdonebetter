// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { UserStatusResponse, APIResponse } from '@dinnerdonebetter/models';

export async function getAuthStatus(client: Axios): Promise<APIResponse<UserStatusResponse>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<UserStatusResponse>>(`/auth/status`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
