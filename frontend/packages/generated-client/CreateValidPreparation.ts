// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidPreparation, APIResponse, ValidPreparationCreationRequestInput } from '@dinnerdonebetter/models';

export async function createValidPreparation(
  client: Axios,

  input: ValidPreparationCreationRequestInput,
): Promise<APIResponse<ValidPreparation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
