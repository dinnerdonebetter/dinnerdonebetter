// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidIngredientGroup, APIResponse, ValidIngredientGroupUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateValidIngredientGroup(
  client: Axios,
  validIngredientGroupID: string,
  input: ValidIngredientGroupUpdateRequestInput,
): Promise<APIResponse<ValidIngredientGroup>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidIngredientGroup>>(
      `/api/v1/valid_ingredient_groups/${validIngredientGroupID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
