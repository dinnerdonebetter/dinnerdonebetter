import { buildServerSideLogger } from '@dinnerdonebetter/logger';
import { buildAPIProxyRoute, UserSessionDetails } from '@dinnerdonebetter/next-routes';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { encryptorDecryptor } from '../../../src/encryption';
import { webappCookieName } from '../../../src/constants';

const logger = buildServerSideLogger('admin_v1_api_proxy');
const apiProxyRoute = buildAPIProxyRoute(
  logger,
  webappCookieName,
  encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>,
);

export default apiProxyRoute;
