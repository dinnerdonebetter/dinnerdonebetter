// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientState, 
  APIResponse, 
  ValidIngredientStateUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateValidIngredientState(
  client: Axios,
  validIngredientStateID: string,
  input: ValidIngredientStateUpdateRequestInput,
): Promise<  APIResponse <  ValidIngredientState >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < ValidIngredientState  >  >(`/api/v1/valid_ingredient_states/${validIngredientStateID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}