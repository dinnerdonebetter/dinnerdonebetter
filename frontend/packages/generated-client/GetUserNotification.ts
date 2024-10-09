// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserNotification, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getUserNotification(
  client: Axios,
  userNotificationID: string,
	): Promise<  APIResponse <  UserNotification >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < UserNotification  >  >(`/api/v1/user_notifications/${userNotificationID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}