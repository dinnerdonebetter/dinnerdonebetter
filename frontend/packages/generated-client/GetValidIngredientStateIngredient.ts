// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidIngredientStateIngredient, APIResponse } from '@dinnerdonebetter/models';

export async function getValidIngredientStateIngredient(
  client: Axios,
  validIngredientStateIngredientID: string,
): Promise<APIResponse<ValidIngredientStateIngredient>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredientStateIngredient>>(
      `/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
