// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { Household, APIResponse, HouseholdUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateHousehold(
  client: Axios,
  householdID: string,
  input: HouseholdUpdateRequestInput,
): Promise<APIResponse<Household>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<Household>>(`/api/v1/households/${householdID}`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
