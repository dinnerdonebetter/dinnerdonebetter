// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserCreationResponse, 
  APIResponse, 
  UserRegistrationInput, 
} from '@dinnerdonebetter/models';

export async function createUser(
  client: Axios,
  
  input: UserRegistrationInput,
): Promise<  APIResponse <  UserCreationResponse >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < UserCreationResponse  >  >(`/users`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}