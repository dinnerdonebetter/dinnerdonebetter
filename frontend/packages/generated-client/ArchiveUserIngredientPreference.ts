// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserIngredientPreference, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveUserIngredientPreference(
  client: Axios,
  userIngredientPreferenceID: string,
	): Promise<  APIResponse <  UserIngredientPreference >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < UserIngredientPreference  >  >(`/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}