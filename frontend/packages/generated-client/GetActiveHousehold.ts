// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { Household, APIResponse } from '@dinnerdonebetter/models';

export async function getActiveHousehold(client: Axios): Promise<APIResponse<Household>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Household>>(`/api/v1/households/current`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
