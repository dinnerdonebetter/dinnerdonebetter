// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { Meal, APIResponse } from '@dinnerdonebetter/models';

export async function archiveMeal(client: Axios, mealID: string): Promise<APIResponse<Meal>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<Meal>>(`/api/v1/meals/${mealID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
