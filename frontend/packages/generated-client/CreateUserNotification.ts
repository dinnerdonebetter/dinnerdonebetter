// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { UserNotification, APIResponse, UserNotificationCreationRequestInput } from '@dinnerdonebetter/models';

export async function createUserNotification(
  client: Axios,

  input: UserNotificationCreationRequestInput,
): Promise<APIResponse<UserNotification>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<UserNotification>>(`/api/v1/user_notifications`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
