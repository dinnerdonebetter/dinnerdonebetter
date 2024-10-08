// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidPreparation, APIResponse, ValidPreparationUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateValidPreparation(
  client: Axios,
  validPreparationID: string,
  input: ValidPreparationUpdateRequestInput,
): Promise<APIResponse<ValidPreparation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidPreparation>>(
      `/api/v1/valid_preparations/${validPreparationID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
