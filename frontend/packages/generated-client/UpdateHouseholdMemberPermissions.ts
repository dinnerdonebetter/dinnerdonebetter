// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { UserPermissionsResponse, APIResponse, ModifyUserPermissionsInput } from '@dinnerdonebetter/models';

export async function updateHouseholdMemberPermissions(
  client: Axios,
  householdID: string,
  userID: string,
  input: ModifyUserPermissionsInput,
): Promise<APIResponse<UserPermissionsResponse>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.patch<APIResponse<UserPermissionsResponse>>(
      `/api/v1/households/${householdID}/members/${userID}/permissions`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
