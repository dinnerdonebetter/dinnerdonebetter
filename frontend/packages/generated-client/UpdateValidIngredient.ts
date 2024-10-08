// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidIngredient, APIResponse, ValidIngredientUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateValidIngredient(
  client: Axios,
  validIngredientID: string,
  input: ValidIngredientUpdateRequestInput,
): Promise<APIResponse<ValidIngredient>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidIngredient>>(
      `/api/v1/valid_ingredients/${validIngredientID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
