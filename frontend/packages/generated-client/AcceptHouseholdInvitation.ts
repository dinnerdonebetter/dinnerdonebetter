// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { HouseholdInvitation, APIResponse, HouseholdInvitationUpdateRequestInput } from '@dinnerdonebetter/models';

export async function acceptHouseholdInvitation(
  client: Axios,
  householdInvitationID: string,
  input: HouseholdInvitationUpdateRequestInput,
): Promise<APIResponse<HouseholdInvitation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<HouseholdInvitation>>(
      `/api/v1/household_invitations/${householdInvitationID}/accept`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
