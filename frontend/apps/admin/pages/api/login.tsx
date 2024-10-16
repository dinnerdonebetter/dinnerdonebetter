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
  oauth2ClientID: process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID || '',
  oauth2ClientSecret: process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET || '',
  serverSideTracer,
  cookieFunc: encodeCookie,
  cookieName: webappCookieName,
  encryptorDecryptor: encryptorDecryptor as EncryptorDecryptor<UserSessionDetails>,
  admin: true,
});
