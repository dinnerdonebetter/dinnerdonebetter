import { buildServerSideClientBuilder } from '@dinnerdonebetter/api-client';

import { apiCookieName } from '../constants';

export const buildServerSideClient = buildServerSideClientBuilder(apiCookieName);
