import { buildServerSideLogger } from '@dinnerdonebetter/logger';
import { buildLogoutRoute, UserSessionDetails } from '@dinnerdonebetter/next-routes';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { webappCookieName } from '../../src/constants';
import { serverSideTracer } from '../../src/tracer';
import { encryptorDecryptor } from '../../src/encryption';

const logger = buildServerSideLogger('logout_route');

const logoutRoute = buildLogoutRoute(
  logger,
  serverSideTracer,
  webappCookieName,
  encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>,
);

export default logoutRoute;
