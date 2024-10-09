// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  User, 
  APIResponse, 
  EmailAddressVerificationRequestInput, 
} from '@dinnerdonebetter/models';

export async function verifyEmailAddress(
  client: Axios,
  
  input: EmailAddressVerificationRequestInput,
): Promise<  APIResponse <  User >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < User  >  >(`/users/email_address/verify`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}