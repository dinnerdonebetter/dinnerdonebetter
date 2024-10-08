// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { User, APIResponse, AvatarUpdateInput } from '@dinnerdonebetter/models';

export async function uploadUserAvatar(
  client: Axios,

  input: AvatarUpdateInput,
): Promise<APIResponse<User>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<User>>(`/api/v1/users/avatar/upload`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
