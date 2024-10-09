// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  PasswordResetResponse, 
  APIResponse, 
  PasswordUpdateInput, 
} from '@dinnerdonebetter/models';

export async function updatePassword(
  client: Axios,
  
  input: PasswordUpdateInput,
): Promise<  APIResponse <  PasswordResetResponse >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < PasswordResetResponse  >  >(`/api/v1/users/password/new`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}