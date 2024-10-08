// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { HouseholdInvitation, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getSentHouseholdInvitations(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<HouseholdInvitation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<HouseholdInvitation>>>(`/api/v1/household_invitations/sent`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<HouseholdInvitation>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
