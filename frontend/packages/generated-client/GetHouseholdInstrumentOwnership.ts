// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { HouseholdInstrumentOwnership, APIResponse } from '@dinnerdonebetter/models';

export async function getHouseholdInstrumentOwnership(
  client: Axios,
  householdInstrumentOwnershipID: string,
): Promise<APIResponse<HouseholdInstrumentOwnership>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<HouseholdInstrumentOwnership>>(
      `/api/v1/households/instruments/${householdInstrumentOwnershipID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
