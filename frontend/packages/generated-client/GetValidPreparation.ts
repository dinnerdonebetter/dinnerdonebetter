// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidPreparation, APIResponse } from '@dinnerdonebetter/models';

export async function getValidPreparation(
  client: Axios,
  validPreparationID: string,
): Promise<APIResponse<ValidPreparation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidPreparation>>(
      `/api/v1/valid_preparations/${validPreparationID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
