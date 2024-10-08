// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserIngredientPreference,
  APIResponse,
  UserIngredientPreferenceUpdateRequestInput,
} from '@dinnerdonebetter/models';

export async function updateUserIngredientPreference(
  client: Axios,
  userIngredientPreferenceID: string,
  input: UserIngredientPreferenceUpdateRequestInput,
): Promise<APIResponse<UserIngredientPreference>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<UserIngredientPreference>>(
      `/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
