import { buildServerSideLogger } from '@dinnerdonebetter/logger';
import { buildAPIProxyRoute } from '@dinnerdonebetter/next-routes';

import { apiCookieName } from '../../../src/constants';

const logger = buildServerSideLogger('webapp_v1_api_proxy');
const apiProxyRoute = buildAPIProxyRoute(logger, apiCookieName);

export default apiProxyRoute;
