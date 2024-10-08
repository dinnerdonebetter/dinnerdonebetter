// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { HouseholdUserMembership, APIResponse } from '@dinnerdonebetter/models';

export async function archiveUserMembership(
  client: Axios,
  householdID: string,
  userID: string,
): Promise<APIResponse<HouseholdUserMembership>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<HouseholdUserMembership>>(
      `/api/v1/households/${householdID}/members/${userID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
