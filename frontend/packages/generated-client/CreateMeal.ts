// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { Meal, APIResponse, MealCreationRequestInput } from '@dinnerdonebetter/models';

export async function createMeal(
  client: Axios,

  input: MealCreationRequestInput,
): Promise<APIResponse<Meal>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<Meal>>(`/api/v1/meals`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
