// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserNotification, 
  APIResponse, 
  UserNotificationUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateUserNotification(
  client: Axios,
  userNotificationID: string,
  input: UserNotificationUpdateRequestInput,
): Promise<  APIResponse <  UserNotification >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.patch<APIResponse < UserNotification  >  >(`/api/v1/user_notifications/${userNotificationID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}