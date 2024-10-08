// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidInstrument, APIResponse } from '@dinnerdonebetter/models';

export async function getRandomValidInstrument(client: Axios): Promise<APIResponse<ValidInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidInstrument>>(`/api/v1/valid_instruments/random`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
