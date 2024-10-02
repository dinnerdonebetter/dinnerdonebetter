import { buildLoginRoute, cookieEncoderBuilder, UserSessionDetails } from '@dinnerdonebetter/next-routes';
import { EncryptorDecryptor } from '@dinnerdonebetter/encryption';

import { encryptorDecryptor } from '../../src/encryption';
import { serverSideTracer } from '../../src/tracer';
import { webappCookieName } from '../../src/constants';

const encodeCookie = cookieEncoderBuilder(
  webappCookieName,
  encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>,
);

export default buildLoginRoute({
  baseURL: 'https://api.dinnerdonebetter.dev',
  scope: 'service_admin',
  oauth2ClientID: '039f605f4f8c13a20dc6639e4095e66d',
  oauth2ClientSecret: 'dc707e2234a05fdfc3d2eb2967c272a2',
  serverSideTracer,
  cookieFunc: encodeCookie,
  admin: true,
});
