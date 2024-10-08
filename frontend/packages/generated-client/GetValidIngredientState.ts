// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidIngredientState, APIResponse } from '@dinnerdonebetter/models';

export async function getValidIngredientState(
  client: Axios,
  validIngredientStateID: string,
): Promise<APIResponse<ValidIngredientState>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredientState>>(
      `/api/v1/valid_ingredient_states/${validIngredientStateID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
