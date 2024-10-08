// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { User, APIResponse, UsernameUpdateInput } from '@dinnerdonebetter/models';

export async function updateUserUsername(
  client: Axios,

  input: UsernameUpdateInput,
): Promise<APIResponse<User>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<User>>(`/api/v1/users/username`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
