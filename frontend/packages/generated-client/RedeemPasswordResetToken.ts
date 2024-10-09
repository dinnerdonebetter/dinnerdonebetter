// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  User, 
  APIResponse, 
  PasswordResetTokenRedemptionRequestInput, 
} from '@dinnerdonebetter/models';

export async function redeemPasswordResetToken(
  client: Axios,
  
  input: PasswordResetTokenRedemptionRequestInput,
): Promise<  APIResponse <  User >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < User  >  >(`/users/password/reset/redeem`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}