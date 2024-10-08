// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { TOTPSecretRefreshResponse, APIResponse, TOTPSecretRefreshInput } from '@dinnerdonebetter/models';

export async function refreshTOTPSecret(
  client: Axios,

  input: TOTPSecretRefreshInput,
): Promise<APIResponse<TOTPSecretRefreshResponse>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<TOTPSecretRefreshResponse>>(`/api/v1/users/totp_secret/new`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
