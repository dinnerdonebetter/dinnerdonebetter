// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientGroup, 
  APIResponse, 
  ValidIngredientGroupCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidIngredientGroup(
  client: Axios,
  
  input: ValidIngredientGroupCreationRequestInput,
): Promise<  APIResponse <  ValidIngredientGroup >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidIngredientGroup  >  >(`/api/v1/valid_ingredient_groups`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}