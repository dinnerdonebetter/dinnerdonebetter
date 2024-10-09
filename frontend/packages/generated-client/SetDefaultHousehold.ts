// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  Household, 
  APIResponse, 
  
} from '@dinnerdonebetter/models';

export async function setDefaultHousehold(
  client: Axios,
  householdID: string,
  
): Promise<  APIResponse <  Household >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < Household  >  >(`/api/v1/households/${householdID}/default`);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}