// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredient, 
  APIResponse, 
  ValidIngredientCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidIngredient(
  client: Axios,
  
  input: ValidIngredientCreationRequestInput,
): Promise<  APIResponse <  ValidIngredient >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidIngredient  >  >(`/api/v1/valid_ingredients`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}