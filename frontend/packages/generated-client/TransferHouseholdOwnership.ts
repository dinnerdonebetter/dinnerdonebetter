// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { Household, APIResponse, HouseholdOwnershipTransferInput } from '@dinnerdonebetter/models';

export async function transferHouseholdOwnership(
  client: Axios,
  householdID: string,
  input: HouseholdOwnershipTransferInput,
): Promise<APIResponse<Household>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<Household>>(`/api/v1/households/${householdID}/transfer`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
