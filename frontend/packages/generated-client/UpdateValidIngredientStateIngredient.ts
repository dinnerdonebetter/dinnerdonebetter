// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientStateIngredient,
  APIResponse,
  ValidIngredientStateIngredientUpdateRequestInput,
} from '@dinnerdonebetter/models';

export async function updateValidIngredientStateIngredient(
  client: Axios,
  validIngredientStateIngredientID: string,
  input: ValidIngredientStateIngredientUpdateRequestInput,
): Promise<APIResponse<ValidIngredientStateIngredient>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidIngredientStateIngredient>>(
      `/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
