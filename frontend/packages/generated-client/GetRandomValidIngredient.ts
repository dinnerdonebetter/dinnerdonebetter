// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidIngredient, APIResponse } from '@dinnerdonebetter/models';

export async function getRandomValidIngredient(client: Axios): Promise<APIResponse<ValidIngredient>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredient>>(`/api/v1/valid_ingredients/random`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
