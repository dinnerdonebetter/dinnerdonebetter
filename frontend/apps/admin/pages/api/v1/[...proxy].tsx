import { buildServerSideLogger } from '@dinnerdonebetter/logger';
import { UserSessionDetails, buildAPIProxyRoute } from '@dinnerdonebetter/next-routes';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { serverSideTracer } from '../../../src/tracer';
import { encryptorDecryptor } from '../../../src/encryption';
import { webappCookieName } from '../../../src/constants';

const logger = buildServerSideLogger('admin_v1_api_proxy');
const apiProxyRoute = buildAPIProxyRoute(
  logger,
  serverSideTracer,
  webappCookieName,
  encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>,
);

export default apiProxyRoute;
