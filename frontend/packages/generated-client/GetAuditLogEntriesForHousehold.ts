// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { AuditLogEntry, APIResponse } from '@dinnerdonebetter/models';

export async function getAuditLogEntriesForHousehold(client: Axios): Promise<APIResponse<AuditLogEntry>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<AuditLogEntry>>(`/api/v1/audit_log_entries/for_household`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
