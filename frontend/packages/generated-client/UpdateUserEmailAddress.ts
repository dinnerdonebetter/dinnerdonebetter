// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  User, 
  APIResponse, 
  UserEmailAddressUpdateInput, 
} from '@dinnerdonebetter/models';

export async function updateUserEmailAddress(
  client: Axios,
  
  input: UserEmailAddressUpdateInput,
): Promise<  APIResponse <  User >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < User  >  >(`/api/v1/users/email_address`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}