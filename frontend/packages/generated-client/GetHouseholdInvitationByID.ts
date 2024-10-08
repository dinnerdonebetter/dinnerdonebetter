// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { HouseholdInvitation, APIResponse } from '@dinnerdonebetter/models';

export async function getHouseholdInvitationByID(
  client: Axios,
  householdID: string,
  householdInvitationID: string,
): Promise<APIResponse<HouseholdInvitation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<HouseholdInvitation>>(
      `/api/v1/households/${householdID}/invitations/${householdInvitationID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
