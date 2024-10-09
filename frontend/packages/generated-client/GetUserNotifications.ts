// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserNotification, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getUserNotifications(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  ): Promise< QueryFilteredResult< UserNotification >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<UserNotification>  >  >(`/api/v1/user_notifications`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<UserNotification>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}