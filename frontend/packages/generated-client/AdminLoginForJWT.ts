// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  JWTResponse, 
  APIResponse, 
  UserLoginInput, 
} from '@dinnerdonebetter/models';

export async function adminLoginForJWT(
  client: Axios,
  
  input: UserLoginInput,
): Promise<  APIResponse <  JWTResponse >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < JWTResponse  >  >(`/users/login/jwt/admin`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}