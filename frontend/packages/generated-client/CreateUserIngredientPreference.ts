// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserIngredientPreference, 
  APIResponse, 
  UserIngredientPreferenceCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createUserIngredientPreference(
  client: Axios,
  
  input: UserIngredientPreferenceCreationRequestInput,
): Promise<  APIResponse <  UserIngredientPreference >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < UserIngredientPreference  >  >(`/api/v1/user_ingredient_preferences`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}