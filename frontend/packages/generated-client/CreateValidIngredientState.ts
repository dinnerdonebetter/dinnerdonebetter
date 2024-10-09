// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientState, 
  APIResponse, 
  ValidIngredientStateCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidIngredientState(
  client: Axios,
  
  input: ValidIngredientStateCreationRequestInput,
): Promise<  APIResponse <  ValidIngredientState >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidIngredientState  >  >(`/api/v1/valid_ingredient_states`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}