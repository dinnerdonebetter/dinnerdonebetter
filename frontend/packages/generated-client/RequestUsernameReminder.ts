// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  User, 
  APIResponse, 
  UsernameReminderRequestInput, 
} from '@dinnerdonebetter/models';

export async function requestUsernameReminder(
  client: Axios,
  
  input: UsernameReminderRequestInput,
): Promise<  APIResponse <  User >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < User  >  >(`/users/username/reminder`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}