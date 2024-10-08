// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidPreparation, APIResponse } from '@dinnerdonebetter/models';

export async function getRandomValidPreparation(client: Axios): Promise<APIResponse<ValidPreparation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations/random`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
