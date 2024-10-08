// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { UserPermissionsResponse, APIResponse, UserPermissionsRequestInput } from '@dinnerdonebetter/models';

export async function checkPermissions(
  client: Axios,

  input: UserPermissionsRequestInput,
): Promise<APIResponse<UserPermissionsResponse>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<UserPermissionsResponse>>(`/api/v1/users/permissions/check`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
