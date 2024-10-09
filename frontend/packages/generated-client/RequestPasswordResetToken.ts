// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  PasswordResetToken, 
  APIResponse, 
  PasswordResetTokenCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function requestPasswordResetToken(
  client: Axios,
  
  input: PasswordResetTokenCreationRequestInput,
): Promise<  APIResponse <  PasswordResetToken >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < PasswordResetToken  >  >(`/users/password/reset`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}