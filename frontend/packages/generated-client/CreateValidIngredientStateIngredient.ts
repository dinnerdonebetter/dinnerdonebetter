// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientStateIngredient, 
  APIResponse, 
  ValidIngredientStateIngredientCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidIngredientStateIngredient(
  client: Axios,
  
  input: ValidIngredientStateIngredientCreationRequestInput,
): Promise<  APIResponse <  ValidIngredientStateIngredient >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidIngredientStateIngredient  >  >(`/api/v1/valid_ingredient_state_ingredients`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}