// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { User, APIResponse, TOTPSecretVerificationInput } from '@dinnerdonebetter/models';

export async function verifyTOTPSecret(
  client: Axios,

  input: TOTPSecretVerificationInput,
): Promise<APIResponse<User>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<User>>(`/users/totp_secret/verify`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
